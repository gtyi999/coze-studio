import argparse
import hashlib
import math
import re
import uuid

from pymilvus import Collection, CollectionSchema, DataType, FieldSchema, connections, utility


DEFAULT_DOCUMENTS = [
    "Milvus is a vector database built for similarity search.",
    "Coze Studio can store embeddings in Milvus for retrieval workflows.",
    "MySQL is better suited for structured relational data than vector search.",
]
CONTENT_MAX_LENGTH = 8192


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Insert your own text into Milvus with a local demo embedder and run a search."
    )
    parser.add_argument("--host", default="127.0.0.1", help="Milvus host")
    parser.add_argument("--port", default="19530", help="Milvus gRPC port")
    parser.add_argument(
        "--collection",
        default=f"coze_demo_{uuid.uuid4().hex[:8]}",
        help="Collection name to use for the demo",
    )
    parser.add_argument(
        "--keep",
        action="store_true",
        help="Keep the collection after the demo finishes",
    )
    parser.add_argument(
        "--overwrite",
        action="store_true",
        help="Drop the existing collection if the target name already exists",
    )
    parser.add_argument(
        "--text",
        action="append",
        default=[],
        help="Document text to insert. Pass this flag multiple times to add multiple documents.",
    )
    parser.add_argument(
        "--input-file",
        help="Path to a UTF-8 text file. Each non-empty line is inserted as one document.",
    )
    parser.add_argument(
        "--query",
        help="Query text to search for. Defaults to the first inserted document.",
    )
    parser.add_argument(
        "--top-k",
        type=int,
        default=3,
        help="How many nearest matches to return",
    )
    parser.add_argument(
        "--dim",
        type=int,
        default=32,
        help="Vector dimension for the local demo embedder",
    )
    return parser.parse_args()


def build_collection(name: str, dim: int) -> Collection:
    fields = [
        FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=False),
        FieldSchema(name="title", dtype=DataType.VARCHAR, max_length=200),
        FieldSchema(name="content", dtype=DataType.VARCHAR, max_length=CONTENT_MAX_LENGTH),
        FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=dim),
    ]
    schema = CollectionSchema(fields=fields, description="Codex Milvus demo collection")
    collection = Collection(name=name, schema=schema)
    collection.create_index(
        field_name="embedding",
        index_params={"index_type": "AUTOINDEX", "metric_type": "COSINE", "params": {}},
    )
    return collection


def load_documents(args: argparse.Namespace) -> list[str]:
    documents: list[str] = []
    documents.extend(text.strip() for text in args.text if text.strip())

    if args.input_file:
        with open(args.input_file, "r", encoding="utf-8-sig") as file:
            documents.extend(line.strip().lstrip("\ufeff") for line in file if line.strip())

    if not documents:
        return DEFAULT_DOCUMENTS.copy()

    return documents


def iter_features(text: str) -> list[str]:
    normalized = " ".join(text.lower().split())
    if not normalized:
        return ["<empty>"]

    features: list[str] = []
    words = re.findall(r"\w+", normalized, flags=re.UNICODE)
    for word in words:
        features.append(f"w:{word}")
    for index in range(len(words) - 1):
        features.append(f"b:{words[index]}_{words[index + 1]}")

    compact = normalized.replace(" ", "")
    for char in compact:
        features.append(f"c:{char}")
    for index in range(len(compact) - 1):
        features.append(f"g:{compact[index:index + 2]}")

    return features or ["<empty>"]


def embed_text(text: str, dim: int) -> list[float]:
    vector = [0.0] * dim
    for feature in iter_features(text):
        digest = hashlib.blake2b(feature.encode("utf-8"), digest_size=8).digest()
        slot = int.from_bytes(digest[:4], "big") % dim
        sign = 1.0 if digest[4] % 2 == 0 else -1.0
        vector[slot] += sign

    norm = math.sqrt(sum(value * value for value in vector))
    if norm == 0:
        return vector

    return [value / norm for value in vector]


def insert_documents(collection: Collection, documents: list[str], dim: int) -> None:
    ids = list(range(1, len(documents) + 1))
    titles = [f"doc_{index}" for index in ids]
    contents = [document[:CONTENT_MAX_LENGTH] for document in documents]
    embeddings = [embed_text(document, dim) for document in documents]
    collection.insert([ids, titles, contents, embeddings])
    collection.flush()


def preview_text(text: str, limit: int = 80) -> str:
    compact = " ".join(text.split())
    if len(compact) <= limit:
        return compact
    return f"{compact[: limit - 3]}..."


def main() -> int:
    args = parse_args()
    alias = "default"
    collection = None
    documents = load_documents(args)
    query_text = (args.query or documents[0]).strip()

    connections.connect(alias=alias, host=args.host, port=args.port)
    print(f"Connected to {args.host}:{args.port}")
    print(f"Server version: {utility.get_server_version()}")
    print("Embedding mode: local hashing demo")

    try:
        if utility.has_collection(args.collection):
            if not args.overwrite:
                print(
                    f"Collection {args.collection!r} already exists. "
                    "Choose another name or rerun with --overwrite."
                )
                return 1
            utility.drop_collection(args.collection)

        collection = build_collection(args.collection, args.dim)
        insert_documents(collection, documents, args.dim)
        print(f"Collection: {args.collection}")
        print(f"Entities: {collection.num_entities}")
        print(f"Inserted documents: {len(documents)}")
        print(f"Query: {query_text}")

        collection.load()
        results = collection.search(
            data=[embed_text(query_text, args.dim)],
            anns_field="embedding",
            param={"metric_type": "COSINE", "params": {}},
            limit=min(args.top_k, len(documents)),
            output_fields=["title", "content"],
        )

        print("Search results:")
        for rank, hit in enumerate(results[0], start=1):
            title = hit.entity.get("title") if hit.entity else ""
            content = hit.entity.get("content") if hit.entity else ""
            print(
                f"{rank}. id={hit.id} score={hit.score:.6f} title={title} "
                f"text={preview_text(content)}"
            )

        return 0
    finally:
        if collection is not None and not args.keep and utility.has_collection(args.collection):
            collection.drop()
            print(f"Dropped collection: {args.collection}")
        connections.disconnect(alias)


if __name__ == "__main__":
    raise SystemExit(main())

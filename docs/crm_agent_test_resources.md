# CRM Agent Test Resources

This repo now includes a Windows-friendly local setup script for CRM agent smoke testing:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\setup\init_crm_agent_test_resources.ps1
```

The script does the following:

1. checks that `http://localhost:8888` is reachable
2. checks that the local MySQL container is running
3. applies `backend/types/ddl/crm_phase1.sql`
4. creates or logs in three CRM agent test users
5. fetches each user's personal `space_id`
6. seeds CRM demo data into `tenant_id == space_id`
7. verifies the CRM dashboard and customer list APIs
8. writes a manifest to `output/crm-agent-test-resources.json`

## Default test users

- `crm-agent-owner@example.com`
- `crm-agent-sales@example.com`
- `crm-agent-analyst@example.com`

Default password:

```text
Passw0rd!123
```

## Example CRM agent questions

- `How many customers do I have now?`
- `Which sales rep has the best performance?`
- `Show me the top five sales reps this quarter.`
- `Which product will sell best next quarter?`
- `Why do you think that product will sell best next quarter?`

## Notes

- Each user gets an isolated personal space, so CRM demo data is seeded separately per account.
- The current CRM NL query backend is still MVP-stage. For supported prompts, the frontend can fall back to built-in example answers.
- The generated manifest is the source of truth for the real `space_id` values on your local machine.

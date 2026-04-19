/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package coze

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	cozehandler "github.com/coze-dev/coze-studio/backend/api/handler/coze"
)

func RegisterManualCRMRoutes(r *server.Hertz) {
	root := r.Group("/", rootMw()...)
	api := root.Group("/api", _apiMw()...)
	crm := api.Group("/crm")

	dashboard := crm.Group("/dashboard")
	dashboard.GET("/overview", cozehandler.GetDashboardOverview)

	customer := crm.Group("/customer")
	customer.GET("/list", cozehandler.ListCustomers)
	customer.GET("/get", cozehandler.GetCustomer)
	customer.GET("/detail", cozehandler.GetCustomer)
	customer.POST("/create", cozehandler.CreateCustomer)
	customer.POST("/update", cozehandler.UpdateCustomer)
	customer.POST("/delete", cozehandler.DeleteCustomer)

	contact := crm.Group("/contact")
	contact.GET("/list", cozehandler.ListContacts)
	contact.GET("/get", cozehandler.GetContact)
	contact.GET("/detail", cozehandler.GetContact)
	contact.POST("/create", cozehandler.CreateContact)
	contact.POST("/update", cozehandler.UpdateContact)
	contact.POST("/delete", cozehandler.DeleteContact)

	opportunity := crm.Group("/opportunity")
	opportunity.GET("/list", cozehandler.ListOpportunities)
	opportunity.GET("/get", cozehandler.GetOpportunity)
	opportunity.GET("/detail", cozehandler.GetOpportunity)
	opportunity.POST("/create", cozehandler.CreateOpportunity)
	opportunity.POST("/update", cozehandler.UpdateOpportunity)
	opportunity.POST("/delete", cozehandler.DeleteOpportunity)

	followRecord := crm.Group("/follow_record")
	followRecord.GET("/list", cozehandler.ListFollowRecords)
	followRecord.GET("/get", cozehandler.GetFollowRecord)
	followRecord.GET("/detail", cozehandler.GetFollowRecord)
	followRecord.POST("/create", cozehandler.CreateFollowRecord)
	followRecord.POST("/update", cozehandler.UpdateFollowRecord)
	followRecord.POST("/delete", cozehandler.DeleteFollowRecord)

	product := crm.Group("/product")
	product.GET("/list", cozehandler.ListProducts)
	product.GET("/get", cozehandler.GetProduct)
	product.GET("/detail", cozehandler.GetProduct)
	product.POST("/create", cozehandler.CreateProduct)
	product.POST("/update", cozehandler.UpdateProduct)
	product.POST("/delete", cozehandler.DeleteProduct)

	salesOrder := crm.Group("/sales_order")
	salesOrder.GET("/list", cozehandler.ListSalesOrders)
	salesOrder.GET("/get", cozehandler.GetSalesOrder)
	salesOrder.GET("/detail", cozehandler.GetSalesOrder)
	salesOrder.POST("/create", cozehandler.CreateSalesOrder)
	salesOrder.POST("/update", cozehandler.UpdateSalesOrder)
	salesOrder.POST("/delete", cozehandler.DeleteSalesOrder)

	query := crm.Group("/query")
	query.POST("/run", cozehandler.RunCRMNLQuery)
	query.GET("/semantic_catalog", cozehandler.GetCRMSemanticCatalog)
	query.GET("/logs", cozehandler.ListCRMQueryLogs)
	query.GET("/forecast", cozehandler.GetCRMForecastResult)
}

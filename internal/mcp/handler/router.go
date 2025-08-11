package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nickborysov/go-multiclient-example/internal/dependency"
	"github.com/nickborysov/go-multiclient-example/internal/model"
	"github.com/samber/lo"
)

type Router struct {
	service dependency.Service
	*mcp.Server
}

func NewRouter(service dependency.Service) *Router {
	server := mcp.NewServer(&mcp.Implementation{Name: "example", Version: "v1.0.0"}, nil)

	r := &Router{
		Server:  server,
		service: service,
	}
	r.RegisterRoutes()

	return r
}

func (r *Router) HTTPHandler() http.Handler {
	// This method is not used in MCP, but we can keep it for consistency
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return r.Server
	}, nil)
}

func (r *Router) RegisterRoutes() {
	mcp.AddTool(r.Server, &mcp.Tool{Name: "example", Description: "example"}, r.HandleExample)
	mcp.AddTool(r.Server, &mcp.Tool{Name: "sendInfo", Description: "Send me Info"}, r.HandleSendInfo)
}

func (r *Router) HandleExample(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[struct{}],
) (*mcp.CallToolResultFor[model.ExampleResponse], error) {
	resp, err := r.service.GetTestResponse()
	if err != nil {
		log.Fatal(err)
	}

	return &mcp.CallToolResultFor[model.ExampleResponse]{
		StructuredContent: resp,
		Content: []mcp.Content{&mcp.TextContent{
			Text:        resp.Message,
			Meta:        nil,
			Annotations: nil,
		}},
	}, nil
}

func (r *Router) HandleSendInfo(
	ctx context.Context,
	cc *mcp.ServerSession,
	params *mcp.CallToolParamsFor[model.SendInfoRequest],
) (*mcp.CallToolResultFor[model.ExampleResponse], error) {
	log.Printf("Received first name %s\n", params.Arguments.FirstName)
	//log.Printf("Received last name %s\n", params.Arguments.LastName)
	//log.Printf("Received email %s\n", params.Arguments.Email)
	//log.Printf("Received phone %s\n", params.Arguments.Phone)
	//log.Printf("Received address %s\n", params.Arguments.Address)
	roots, err := cc.ListRoots(ctx, nil)
	if err != nil {
		log.Printf("Error listing roots: %v\n", err)
	}

	var rootNames []string
	if roots != nil {
		rootNames = lo.Map(roots.Roots, func(r *mcp.Root, _ int) string {
			return r.Name
		})
	}

	log.Printf("Received session ID %s, roots: %s\n", cc.ID(), strings.Join(rootNames, ","))
	log.Printf("------------------------------\n")
	// Here you can process the received information as needed
	resp, err := r.service.GetTestResponse()
	if err != nil {
		log.Fatal(err)
	}

	return &mcp.CallToolResultFor[model.ExampleResponse]{
		StructuredContent: resp,
		Content: []mcp.Content{&mcp.TextContent{
			Text:        resp.Message,
			Meta:        nil,
			Annotations: nil,
		}},
	}, nil

}

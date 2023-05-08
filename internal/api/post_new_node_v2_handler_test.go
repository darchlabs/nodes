package api

import (
	"testing"
	"time"

	"github.com/darchlabs/nodes/internal/manager"
	"github.com/jaekwon/testify/require"
)

func Test_PostNewNodeV2Handeler_invoke(t *testing.T) {
	testCases := []struct {
		name                    string
		mockCtx                 *Context
		mockInstanceInsertQuery *mockInstanceInsertQuery
		mockNodeManager         mockNodeManager
		inputReq                *postNewNodev2HandlerRequest
		expectedStatus          int
		expectedErrMsg          string
	}{
		{
			name:                    "should not return error",
			mockInstanceInsertQuery: &mockInstanceInsertQuery{err: nil},
			mockCtx: &Context{
				nodeManager: &mockNodeManager{
					res: &manager.NodeInstance{
						ID:        "test-id",
						Name:      "darch-node",
						Artifacts: &manager.Artifacts{Deployments: []string{"darch", "node"}},
						Config: &manager.NodeConfig{
							Network:     "darchlabs",
							Environment: "mainnet",
							Port:        6969,
							CreatedAt:   time.Now(),
						},
					},
				},
			},
			inputReq:       &postNewNodev2HandlerRequest{},
			expectedStatus: 201,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := PostNewNodeV2Handler{
				instanceInsertQuery: tc.mockInstanceInsertQuery.Insert,
			}

			payload, status, err := h.invoke(tc.mockCtx, tc.inputReq)

			require.NoError(t, err)
			if tc.expectedErrMsg != "" {
				require.Equal(t, err.Error(), tc.expectedErrMsg)
				return
			}

			require.Equal(t, status, tc.expectedStatus)
			require.NotNil(t, payload)
		})
	}
}

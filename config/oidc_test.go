package config_test

import (
	"testing"

	"github.com/KMConner/kubectl-auth0/config"
	"github.com/google/go-cmp/cmp"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestLoadOidcConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		k8sConfig api.Config
		cmdArgs   config.CmdLine
		want      config.OidcRequest
		wantName  string
		wantErr   bool
	}{
		{
			name: "Idp url, client id and new username is specified in the command line arguments",
			cmdArgs: config.CmdLine{
				IdpUrl:   "https://example.com/",
				ClientId: "client-1234",
				NewUsername: "new-user1",
			},
			want: config.OidcRequest{
				ClientId: "client-1234",
				IdpUrl:   "https://example.com/",
			},
			wantName: "new-user1",
			wantErr:  false,
		},
		{
			name: "Context specified",
			cmdArgs: config.CmdLine{
				ContextName: "ctx1",
			},
			want: config.OidcRequest{
				ClientId: "client-1234",
				IdpUrl:   "https://example.com/",
			},
			k8sConfig: api.Config{
				AuthInfos: map[string]*api.AuthInfo{
					"user1": {
						AuthProvider: &api.AuthProviderConfig{
							Name: "oidc",
							Config: map[string]string{
								"client-id":      "client-1234",
								"idp-issuer-url": "https://example.com/",
							},
						},
					},
				},
				Contexts: map[string]*api.Context{
					"ctx1": {
						AuthInfo: "user1",
					},
				},
			},
			wantName: "user1",
			wantErr:  false,
		},
		{
			name: "Context specified with new username",
			cmdArgs: config.CmdLine{
				ContextName: "ctx1",
				NewUsername: "new-user1",
			},
			want: config.OidcRequest{
				ClientId: "client-1234",
				IdpUrl:   "https://example.com/",
			},
			k8sConfig: api.Config{
				AuthInfos: map[string]*api.AuthInfo{
					"user1": {
						AuthProvider: &api.AuthProviderConfig{
							Name: "oidc",
							Config: map[string]string{
								"client-id":      "client-1234",
								"idp-issuer-url": "https://example.com/",
							},
						},
					},
				},
				Contexts: map[string]*api.Context{
					"ctx1": {
						AuthInfo: "user1",
					},
				},
			},
			wantName: "new-user1",
			wantErr:  false,
		},

		{
			name: "Use default context",
			want: config.OidcRequest{
				ClientId: "client-1234",
				IdpUrl:   "https://example.com/",
			},
			k8sConfig: api.Config{
				AuthInfos: map[string]*api.AuthInfo{
					"user1": {
						AuthProvider: &api.AuthProviderConfig{
							Name: "oidc",
							Config: map[string]string{
								"client-id":      "client-1234",
								"idp-issuer-url": "https://example.com/",
							},
						},
					},
				},
				Contexts: map[string]*api.Context{
					"ctx1": {
						AuthInfo: "user1",
					},
				},
				CurrentContext: "ctx1",
			},
			wantName: "user1",
			wantErr:  false,
		},
		{
			name:      "Context is not specified",
			k8sConfig: api.Config{},
			wantName:  "user1",
			wantErr:   true,
		},
		{
			name: "Context not found",
			k8sConfig: api.Config{
				AuthInfos: map[string]*api.AuthInfo{
					"user1": {
						AuthProvider: &api.AuthProviderConfig{
							Name: "oidc",
							Config: map[string]string{
								"client-id":      "client-1234",
								"idp-issuer-url": "https://example.com/",
							},
						},
					},
				},
				Contexts: map[string]*api.Context{
					"ctx1": {
						AuthInfo: "user1",
					},
				},
				CurrentContext: "ctx2",
			},
			wantErr: true,
		},
		{
			name: "auth info not found",
			k8sConfig: api.Config{
				AuthInfos: map[string]*api.AuthInfo{
					"user1": {
						AuthProvider: &api.AuthProviderConfig{
							Name: "oidc",
							Config: map[string]string{
								"client-id":      "client-1234",
								"idp-issuer-url": "https://example.com/",
							},
						},
					},
				},
				Contexts: map[string]*api.Context{
					"ctx1": {
						AuthInfo: "user2",
					},
				},
				CurrentContext: "ctx1",
			},
			wantErr: true,
		},
		{
			name: "Different auth provider",
			k8sConfig: api.Config{
				AuthInfos: map[string]*api.AuthInfo{
					"user1": {
						AuthProvider: &api.AuthProviderConfig{
							Name: "auth provider",
						},
					},
				},
				Contexts: map[string]*api.Context{
					"ctx1": {
						AuthInfo: "user1",
					},
				},
				CurrentContext: "ctx1",
			},
			wantErr: true,
		},
		{
			name: "client id is not specified in auth provider config",
			cmdArgs: config.CmdLine{
				ContextName: "ctx1",
			},
			k8sConfig: api.Config{
				AuthInfos: map[string]*api.AuthInfo{
					"user1": {
						AuthProvider: &api.AuthProviderConfig{
							Name: "oidc",
							Config: map[string]string{
								"idp-issuer-url": "https://example.com/",
							},
						},
					},
				},
				Contexts: map[string]*api.Context{
					"ctx1": {
						AuthInfo: "user1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "idp url is not specified in auth provider config",
			cmdArgs: config.CmdLine{
				ContextName: "ctx1",
			},
			k8sConfig: api.Config{
				AuthInfos: map[string]*api.AuthInfo{
					"user1": {
						AuthProvider: &api.AuthProviderConfig{
							Name: "oidc",
							Config: map[string]string{
								"client-id": "client-1234",
							},
						},
					},
				},
				Contexts: map[string]*api.Context{
					"ctx1": {
						AuthInfo: "user1",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase := testCase
			t.Parallel()
			got, gotName, err := config.LoadOidcConfig(&testCase.cmdArgs, &testCase.k8sConfig)
			if err != nil {
				if !testCase.wantErr {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}
			if err == nil && testCase.wantErr {
				t.Errorf("Error epected")
			}
			if *got != testCase.want {
				t.Errorf("LoadAndValidate() got := %v, want := %v", *got, testCase.want)
			}
			if gotName != testCase.wantName {
				t.Errorf("LoadAndValidate() gotName := %v, wantName := %v", gotName, testCase.wantName)
			}
		})
	}
}

func TestToAuthInfo(t *testing.T) {
	t.Run("Convert into auth info", func(t *testing.T) {
		t.Parallel()
		want := api.AuthInfo{
			AuthProvider: &api.AuthProviderConfig{
				Name: "oidc",
				Config: map[string]string{
					"client-id":      "client-1234",
					"id-token":       "token-1234",
					"idp-issuer-url": "url-1234",
					"refresh-token":  "refresh-1234",
				},
			},
		}

		oidcConfig := &config.Oidc{
			ClientId:     "client-1234",
			Token:        "token-1234",
			IdpUrl:       "url-1234",
			RefreshToken: "refresh-1234",
		}
		got := oidcConfig.ToAuthInfo()
		if !cmp.Equal(want, *got) {
			t.Fatalf("got := %v, want := %v", *got, want)
		}
	})
}

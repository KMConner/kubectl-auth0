package config_test

import (
	"os"
	"testing"

	"github.com/KMConner/kubectl-auth0/config"
)

func TestLoadCmdArgs(t *testing.T) {
	t.Parallel()
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name    string
		args    []string
		want    config.CmdLine
		wantErr bool
	}{
		{
			name: "Specify all flags",
			args: []string{
				"--context",
				"name1",
				"--client-id",
				"client-id",
				"--idp-issuer-url",
				"idp-url",
				"--new-username",
				"auth0-user",
			},
			want: config.CmdLine{
				ContextName: "name1",
				ClientId:    "client-id",
				IdpUrl:      "idp-url",
				NewUsername: "auth0-user",
			},
			wantErr: false,
		},
		{
			name: "Specify --context flag only",
			args: []string{
				"--context",
				"name1",
			},
			want: config.CmdLine{
				ContextName: "name1",
			},
			wantErr: false,
		},
		{
			name: "Specify wrong flag",
			args: []string{
				"--foo",
				"bar",
			},
			wantErr: true,
		},
		{
			name: "Specify wrong flag 2",
			args: []string{
				"--help",
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			got, err := config.LoadCmdArgs(testCase.args, devNull)
			if err != nil {
				if !testCase.wantErr {
					t.Fatalf("Undexpected err %v.", err)
				}
			} else {
				if testCase.wantErr {
					t.Fatal("Error expected")
				}
			}
			if testCase.want != *got {
				t.Fatalf("Want := %v, Got := %v", testCase.want, *got)
			}
		})
	}
}

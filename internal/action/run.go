package action

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/go-tfe"
	"github.com/sethvargo/go-githubactions"
)

// OutputString takes an interface and returns a stringified version
func OutputString(i interface{}) (string, error) {
	if i == nil {
		return "", nil
	}

	t := reflect.TypeOf(i)
	switch t.Kind() {
	case
		reflect.Array,
		reflect.Map,
		reflect.Slice,
		reflect.Struct:
		b, err := json.Marshal(i)
		if err != nil {
			return "", err
		}
		return string(b), nil

	case
		reflect.Chan,
		reflect.Func:
		return "", nil

	default:
		return fmt.Sprint(i), nil
	}
}

func Run(inputs Inputs) error {
	ctx := context.Background()

	client, err := tfe.NewClient(&tfe.Config{
		Token:   inputs.Token,
		Address: inputs.Address,
	})
	if err != nil {
		return fmt.Errorf("failed to configure client: %w", err)
	}

	workspace, err := client.Workspaces.Read(ctx, inputs.Organization, inputs.Workspace)
	if err != nil {
		return fmt.Errorf("failed to find workspace %s/%s: %w", inputs.Organization, inputs.Workspace, err)
	}

	stateVersion, err := client.StateVersions.CurrentWithOptions(ctx, workspace.ID, &tfe.StateVersionCurrentOptions{
		Include: "outputs",
	})
	if err != nil {
		return fmt.Errorf("failed to fetch state outputs for workspace %s/%s: %w", inputs.Organization, inputs.Workspace, err)
	}

	for _, o := range stateVersion.Outputs {
		str, err := OutputString(o.Value)
		if err != nil {
			return err
		}

		log.Println(str)

		if o.Sensitive {
			githubactions.AddMask(str)
		}

		githubactions.SetOutput(o.Name, str)
	}

	return nil
}

type Inputs struct {
	Token        string
	Address      string
	Workspace    string
	Organization string
}
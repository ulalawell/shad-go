//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var errUpdApi error
	var errSetApi error
	var errAuth *kvapi.AuthError
	var errConflict *kvapi.ConflictError
	var updatedValue *string = nil

	for {
		getReq := kvapi.GetRequest{Key: key}
		getResp := &kvapi.GetResponse{}

		if !errors.Is(errSetApi, kvapi.ErrKeyNotFound) {

		getReqLoop:
			for {
				getResp, errUpdApi = c.Get(&getReq)

				switch {
				case errUpdApi == nil:
					updatedValue = &getResp.Value
					break getReqLoop
				case errors.Is(errUpdApi, kvapi.ErrKeyNotFound):
					break getReqLoop
				case errors.As(errUpdApi, &errAuth):
					return errUpdApi
				default:
					continue
				}
			}

		} else {
			updatedValue = nil
		}

		var newValue string
		var errUpd error
		var oldVersion uuid.UUID

		newValue, errUpd = updateFn(updatedValue)
		if errUpd != nil {
			return errUpd
		}

		if updatedValue == nil || errors.Is(errSetApi, kvapi.ErrKeyNotFound) {
			oldVersion = uuid.UUID{}
		} else {
			oldVersion = getResp.Version
		}

		rsp := &kvapi.SetRequest{
			Key:        key,
			Value:      newValue,
			OldVersion: oldVersion,
			NewVersion: uuid.Must(uuid.NewV4()),
		}

	setReqLoop:
		for {
			_, errSetApi = c.Set(rsp)

			switch {
			case errors.As(errSetApi, &errAuth):
				return errSetApi
			case errSetApi == nil:
				return nil
			case errors.Is(errSetApi, kvapi.ErrKeyNotFound):
				break setReqLoop
			case errors.As(errSetApi, &errConflict):
				if errConflict.ExpectedVersion == rsp.NewVersion {
					return nil
				}
				break setReqLoop
			default:
				continue
			}

		}

	}

}

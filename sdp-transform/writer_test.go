package sdptransform

import (
	"testing"
)

func TestWrite(t *testing.T) {

	session, err := Parse(sdpStr)
	if err != nil {
		t.Error(err)
	}

	_, err = Write(session)

	if err != nil {
		t.Error(err)
	}

}

package mixins

import (
	"testing"

	"github.com/lookupearth/restful"
)

type MixinsCase struct {
	*GetMethod
	*ListMethod
	*PostMethod
	*PatchMethod
	*PutMethod
	*DeleteMethod
}

func TestMixins(t *testing.T) {
	var resource interface{} = &MixinsCase{
		GetMethod:    &GetMethod{},
		ListMethod:   &ListMethod{},
		PostMethod:   &PostMethod{},
		PatchMethod:  &PatchMethod{},
		PutMethod:    &PutMethod{},
		DeleteMethod: &DeleteMethod{},
	}

	if _, ok := resource.(restful.IGet); !ok {
		t.Errorf("MixinsCase not implement restful.IGet")
	}
	if _, ok := resource.(restful.IList); !ok {
		t.Errorf("MixinsCase not implement restful.IList")
	}
	if _, ok := resource.(restful.IPost); !ok {
		t.Errorf("MixinsCase not implement restful.IPost")
	}
	if _, ok := resource.(restful.IPut); !ok {
		t.Errorf("MixinsCase not implement restful.IPut")
	}
	if _, ok := resource.(restful.IPatch); !ok {
		t.Errorf("MixinsCase not implement restful.IPatch")
	}
	if _, ok := resource.(restful.IDelete); !ok {
		t.Errorf("MixinsCase not implement restful.IDelete")
	}

	if _, ok := resource.(restful.IGetInit); !ok {
		t.Errorf("MixinsCase not implement restful.IGetInit")
	}
	if _, ok := resource.(restful.IListInit); !ok {
		t.Errorf("MixinsCase not implement restful.IListInit")
	}
	if _, ok := resource.(restful.IPostInit); !ok {
		t.Errorf("MixinsCase not implement restful.IPostInit")
	}
	if _, ok := resource.(restful.IPutInit); !ok {
		t.Errorf("MixinsCase not implement restful.IPutInit")
	}
	if _, ok := resource.(restful.IPatchInit); !ok {
		t.Errorf("MixinsCase not implement restful.IPatchInit")
	}
	if _, ok := resource.(restful.IDeleteInit); !ok {
		t.Errorf("MixinsCase not implement restful.IDeleteInit")
	}

}

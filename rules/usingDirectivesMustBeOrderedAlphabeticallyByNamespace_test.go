package rules

import (
	"bytes"
	"fmt"
	"testing"
)

func TestUsingDirectivesMustBeOrderedAlphabeticallyByNamespace(t *testing.T) {
	input := []byte(`using Company;
using CompanyB.System;
using Company.Collections.Generic;
using Company.Linq;

namespace Company.Blah {}
using (var something = new Something()) {}`)
	expected := []byte(`using Company;
using Company.Collections.Generic;
using Company.Linq;
using CompanyB.System;

namespace Company.Blah {}
using (var something = new Something()) {}`)

	actual := applyUsingDirectivesMustBeOrderedAlphabeticallyByNamespace(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}

	input = []byte(`using System;
using System.Collections.Generic;
using System.Linq;
using System.Web.Services;
using System.Web.UI;
using Usar.Eks.ProjDoc.Concrete;
using Usar.Eks.ProjDoc.Extensions;

namespace Company.Blah {}`)
	expected = []byte(`using System;
using System.Collections.Generic;
using System.Linq;
using System.Web.Services;
using System.Web.UI;
using Usar.Eks.ProjDoc.Concrete;
using Usar.Eks.ProjDoc.Extensions;

namespace Company.Blah {}`)

	actual = applyUsingDirectivesMustBeOrderedAlphabeticallyByNamespace(input)
	if !bytes.Equal(expected, actual) {
		fmt.Println(string(actual))
		t.Fail()
	}
}

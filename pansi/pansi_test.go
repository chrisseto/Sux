package pansi

import "testing"

func TestInvalidStrings(t *testing.T) {
	invalid := []string{"", "", " ", "OOOOOO", "[", "jj"}

	for _, s := range invalid {
		result := Parse(s)
		if result != nil {
			t.Errorf("AnsiEscapeCode %v is not nil", result.Type)
		}
	}
}

func TestEscapeState(t *testing.T) {
	p := NewParser()

	if p.Result() != nil {
		t.Errorf("NewParser's result is not nil")
	}

	p.Feed('')

	_, ok := p.state.(*escape)

	if !ok {
		t.Errorf("%T is not %T; %+v", p.state, escape{}, ok)
	}

}

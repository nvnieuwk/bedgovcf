package bedgovcf

import "testing"

func TestFieldResolving(t *testing.T) {
	value, _ := resolveField([]string{"$test"}, []string{"value"}, []string{"test"})
	if value != "value" {
		t.Fatalf("Expected value to be 'value', got %s", value)
	}

	value, _ = resolveField([]string{"$test", "$test2"}, []string{"value", "I don't want this", "value2"}, []string{"test", "whut", "test2"})
	if value != "value value2" {
		t.Fatalf("Expected value to be 'value value2', got %s", value)
	}

	value, _ = resolveField([]string{"~sum", "$test", "$test2"}, []string{"10", "I don't want this", "2"}, []string{"test", "whut", "test2"})
	if value != "12" {
		t.Fatalf("Expected value to be '12', got %s", value)
	}
}

func TestRound(t *testing.T) {
	value, _ := resolveField([]string{"~round", "1.5"}, []string{}, []string{})
	if value != "2" {
		t.Fatalf("Expected value to be '2', got %s", value)
	}

	value, _ = resolveField([]string{"~round", "1.4"}, []string{}, []string{})
	if value != "1" {
		t.Fatalf("Expected value to be '1', got %s", value)
	}

	value, _ = resolveField([]string{"~round", "-2.3"}, []string{}, []string{})
	if value != "-2" {
		t.Fatalf("Expected value to be '-2', got %s", value)
	}

	value, _ = resolveField([]string{"~round", "-0.2"}, []string{}, []string{})
	if value != "0" {
		t.Fatalf("Expected value to be '0', got %s", value)
	}

	value, _ = resolveField([]string{"~round", "45698742.2"}, []string{}, []string{})
	if value != "45698742" {
		t.Fatalf("Expected value to be '45698742', got %s", value)
	}
}

func TestSum(t *testing.T) {
	value, _ := resolveField([]string{"~sum", "1.5"}, []string{}, []string{})
	if value != "1.5" {
		t.Fatalf("Expected value to be '1.5', got %s", value)
	}

	value, _ = resolveField([]string{"~sum", "1", "2"}, []string{}, []string{})
	if value != "3" {
		t.Fatalf("Expected value to be '3', got %s", value)
	}

	value, _ = resolveField([]string{"~sum", "-2", "2"}, []string{}, []string{})
	if value != "0" {
		t.Fatalf("Expected value to be '0', got %s", value)
	}

	value, _ = resolveField([]string{"~sum", "10", "15", "20", "-5"}, []string{}, []string{})
	if value != "40" {
		t.Fatalf("Expected value to be '40', got %s", value)
	}

	value, _ = resolveField([]string{"~sum", "-10", "-20", "5"}, []string{}, []string{})
	if value != "-25" {
		t.Fatalf("Expected value to be '-25', got %s", value)
	}

	value, _ = resolveField([]string{"~sum", "12500000", "2500000"}, []string{}, []string{})
	if value != "15000000" {
		t.Fatalf("Expected value to be '15000000', got %s", value)
	}
}

func TestMin(t *testing.T) {
	value, _ := resolveField([]string{"~min", "1.5"}, []string{}, []string{})
	if value != "1.5" {
		t.Fatalf("Expected value to be '1.5', got %s", value)
	}

	value, _ = resolveField([]string{"~min", "1", "2"}, []string{}, []string{})
	if value != "-1" {
		t.Fatalf("Expected value to be '-1', got %s", value)
	}

	value, _ = resolveField([]string{"~min", "-2", "2"}, []string{}, []string{})
	if value != "-4" {
		t.Fatalf("Expected value to be '-4', got %s", value)
	}

	value, _ = resolveField([]string{"~min", "50", "15", "20", "-5"}, []string{}, []string{})
	if value != "20" {
		t.Fatalf("Expected value to be '20', got %s", value)
	}

	value, _ = resolveField([]string{"~min", "-10", "-20", "5"}, []string{}, []string{})
	if value != "5" {
		t.Fatalf("Expected value to be '5', got %s", value)
	}

	value, _ = resolveField([]string{"~min", "12500000", "2500000"}, []string{}, []string{})
	if value != "10000000" {
		t.Fatalf("Expected value to be '10000000', got %s", value)
	}
}

func TestIf(t *testing.T) {
	value, _ := resolveField([]string{"~if", "1", ">", "2", "true", "false"}, []string{}, []string{})
	if value != "false" {
		t.Fatalf("Expected value to be 'false', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "1", "<", "2", "true", "false"}, []string{}, []string{})
	if value != "true" {
		t.Fatalf("Expected value to be 'true', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "1", "<=", "2", "true", "false"}, []string{}, []string{})
	if value != "true" {
		t.Fatalf("Expected value to be 'true', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "2", "<=", "2", "true", "false"}, []string{}, []string{})
	if value != "true" {
		t.Fatalf("Expected value to be 'true', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "1", ">=", "2", "true", "false"}, []string{}, []string{})
	if value != "false" {
		t.Fatalf("Expected value to be 'false', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "2", ">=", "2", "true", "false"}, []string{}, []string{})
	if value != "true" {
		t.Fatalf("Expected value to be 'true', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "1", "==", "2", "true", "false"}, []string{}, []string{})
	if value != "false" {
		t.Fatalf("Expected value to be 'false', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "1", "==", "1", "true", "false"}, []string{}, []string{})
	if value != "true" {
		t.Fatalf("Expected value to be 'true', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "test", "==", "test2", "true", "false"}, []string{}, []string{})
	if value != "false" {
		t.Fatalf("Expected value to be 'false', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "test", "==", "test", "true", "false"}, []string{}, []string{})
	if value != "true" {
		t.Fatalf("Expected value to be 'true', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "test", "==", "test2", "true", "~sum", "10", "20"}, []string{}, []string{})
	if value != "30" {
		t.Fatalf("Expected value to be '30', got %s", value)
	}

	value, _ = resolveField([]string{"~if", "test", "==", "test", "true", "~sum", "10", "20"}, []string{}, []string{})
	if value != "true" {
		t.Fatalf("Expected value to be 'true', got %s", value)
	}
}

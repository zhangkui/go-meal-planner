package service

import "testing"

func TestConsumeExactInventoryLeavesZero(t *testing.T) {
	_, _, inventory, _, _ := testServices()
	if _, err := inventory.Add("Flour", "g", 500); err != nil {
		t.Fatal(err)
	}
	item, err := inventory.Consume("Flour", 500)
	if err != nil {
		t.Fatalf("exact consumption failed: %v", err)
	}
	if item.Quantity != 0 {
		t.Fatalf("quantity = %v, want 0", item.Quantity)
	}
}

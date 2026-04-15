package output

import (
	"encoding/json"
	"testing"
)

// TestAnyOfInlineCatType verifies the Cat type fields are accessible.
func TestAnyOfInlineCatType(t *testing.T) {
	id := "cat-1"
	name := "Whiskers"
	breed := "Siamese"
	color := "cream"
	purrs := true

	cat := Cat{
		ID:    &id,
		Name:  &name,
		Breed: &breed,
		Color: &color,
		Purrs: &purrs,
	}

	if *cat.ID != "cat-1" {
		t.Errorf("ID = %q, want %q", *cat.ID, "cat-1")
	}
	if *cat.Name != "Whiskers" {
		t.Errorf("Name = %q, want %q", *cat.Name, "Whiskers")
	}
	if *cat.Breed != "Siamese" {
		t.Errorf("Breed = %q, want %q", *cat.Breed, "Siamese")
	}
	if *cat.Color != "cream" {
		t.Errorf("Color = %q, want %q", *cat.Color, "cream")
	}
	if *cat.Purrs != true {
		t.Errorf("Purrs = %v, want true", *cat.Purrs)
	}
}

// TestAnyOfInlineDogType verifies the Dog type fields are accessible.
func TestAnyOfInlineDogType(t *testing.T) {
	id := "dog-1"
	name := "Rex"
	breed := "Labrador"
	color := "golden"
	barks := true

	dog := Dog{
		ID:    &id,
		Name:  &name,
		Breed: &breed,
		Color: &color,
		Barks: &barks,
	}

	if *dog.ID != "dog-1" {
		t.Errorf("ID = %q, want %q", *dog.ID, "dog-1")
	}
	if *dog.Name != "Rex" {
		t.Errorf("Name = %q, want %q", *dog.Name, "Rex")
	}
	if *dog.Barks != true {
		t.Errorf("Barks = %v, want true", *dog.Barks)
	}
}

// TestAnyOfInlineRatType verifies the Rat type fields are accessible.
func TestAnyOfInlineRatType(t *testing.T) {
	id := "rat-1"
	name := "Remy"
	color := "grey"
	squeaks := true

	rat := Rat{
		ID:      &id,
		Name:    &name,
		Color:   &color,
		Squeaks: &squeaks,
	}

	if *rat.ID != "rat-1" {
		t.Errorf("ID = %q, want %q", *rat.ID, "rat-1")
	}
	if *rat.Name != "Remy" {
		t.Errorf("Name = %q, want %q", *rat.Name, "Remy")
	}
	if *rat.Squeaks != true {
		t.Errorf("Squeaks = %v, want true", *rat.Squeaks)
	}
}

// TestAnyOfInlineFromCatRoundTrip verifies FromCat -> MarshalJSON ->
// UnmarshalJSON -> AsCat round-trip.
func TestAnyOfInlineFromCatRoundTrip(t *testing.T) {
	id := "cat-1"
	name := "Whiskers"
	breed := "Siamese"
	color := "cream"
	purrs := true

	var union GetPets200ResponseJSON2
	if err := union.FromCat(Cat{
		ID:    &id,
		Name:  &name,
		Breed: &breed,
		Color: &color,
		Purrs: &purrs,
	}); err != nil {
		t.Fatalf("FromCat failed: %v", err)
	}

	data, err := union.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var decoded GetPets200ResponseJSON2
	if err := decoded.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	got, err := decoded.AsCat()
	if err != nil {
		t.Fatalf("AsCat failed: %v", err)
	}
	if *got.ID != "cat-1" {
		t.Errorf("ID = %q, want %q", *got.ID, "cat-1")
	}
	if *got.Name != "Whiskers" {
		t.Errorf("Name = %q, want %q", *got.Name, "Whiskers")
	}
	if *got.Breed != "Siamese" {
		t.Errorf("Breed = %q, want %q", *got.Breed, "Siamese")
	}
	if *got.Purrs != true {
		t.Errorf("Purrs = %v, want true", *got.Purrs)
	}
}

// TestAnyOfInlineFromDogRoundTrip verifies FromDog -> MarshalJSON ->
// UnmarshalJSON -> AsDog round-trip.
func TestAnyOfInlineFromDogRoundTrip(t *testing.T) {
	id := "dog-1"
	name := "Buddy"
	barks := true

	var union GetPets200ResponseJSON2
	if err := union.FromDog(Dog{
		ID:    &id,
		Name:  &name,
		Barks: &barks,
	}); err != nil {
		t.Fatalf("FromDog failed: %v", err)
	}

	data, err := union.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var decoded GetPets200ResponseJSON2
	if err := decoded.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	got, err := decoded.AsDog()
	if err != nil {
		t.Fatalf("AsDog failed: %v", err)
	}
	if *got.ID != "dog-1" {
		t.Errorf("ID = %q, want %q", *got.ID, "dog-1")
	}
	if *got.Name != "Buddy" {
		t.Errorf("Name = %q, want %q", *got.Name, "Buddy")
	}
	if *got.Barks != true {
		t.Errorf("Barks = %v, want true", *got.Barks)
	}
}

// TestAnyOfInlineFromRatRoundTrip verifies FromRat -> MarshalJSON ->
// UnmarshalJSON -> AsRat round-trip.
func TestAnyOfInlineFromRatRoundTrip(t *testing.T) {
	id := "rat-1"
	name := "Remy"
	squeaks := true

	var union GetPets200ResponseJSON2
	if err := union.FromRat(Rat{
		ID:      &id,
		Name:    &name,
		Squeaks: &squeaks,
	}); err != nil {
		t.Fatalf("FromRat failed: %v", err)
	}

	data, err := union.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var decoded GetPets200ResponseJSON2
	if err := decoded.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	got, err := decoded.AsRat()
	if err != nil {
		t.Fatalf("AsRat failed: %v", err)
	}
	if *got.ID != "rat-1" {
		t.Errorf("ID = %q, want %q", *got.ID, "rat-1")
	}
	if *got.Name != "Remy" {
		t.Errorf("Name = %q, want %q", *got.Name, "Remy")
	}
	if *got.Squeaks != true {
		t.Errorf("Squeaks = %v, want true", *got.Squeaks)
	}
}

// TestAnyOfInlineUnmarshalJSONObject verifies that raw JSON can be unmarshaled
// into the union and then extracted as the correct variant.
func TestAnyOfInlineUnmarshalJSONObject(t *testing.T) {
	input := `{"id":"pet-1","name":"Furball","color":"brown","purrs":true}`

	var union GetPets200ResponseJSON2
	if err := union.UnmarshalJSON([]byte(input)); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// The JSON has "purrs", so it should decode as a Cat.
	cat, err := union.AsCat()
	if err != nil {
		t.Fatalf("AsCat failed: %v", err)
	}
	if *cat.Name != "Furball" {
		t.Errorf("Cat.Name = %q, want %q", *cat.Name, "Furball")
	}
	if *cat.Color != "brown" {
		t.Errorf("Cat.Color = %q, want %q", *cat.Color, "brown")
	}
	if *cat.Purrs != true {
		t.Errorf("Cat.Purrs = %v, want true", *cat.Purrs)
	}

	// anyOf allows the same data to also be read as Dog or Rat (shared fields
	// decode, variant-specific fields are zero/nil).
	dog, err := union.AsDog()
	if err != nil {
		t.Fatalf("AsDog failed: %v", err)
	}
	if *dog.Name != "Furball" {
		t.Errorf("Dog.Name = %q, want %q", *dog.Name, "Furball")
	}
	if dog.Barks != nil {
		t.Errorf("Dog.Barks = %v, want nil (not in input)", *dog.Barks)
	}

	rat, err := union.AsRat()
	if err != nil {
		t.Fatalf("AsRat failed: %v", err)
	}
	if *rat.Name != "Furball" {
		t.Errorf("Rat.Name = %q, want %q", *rat.Name, "Furball")
	}
	if rat.Squeaks != nil {
		t.Errorf("Rat.Squeaks = %v, want nil (not in input)", *rat.Squeaks)
	}
}

// TestAnyOfInlineResponseJSONRoundTrip verifies JSON round-trip for the
// GetPetsJSONResponse wrapper containing anyOf union items.
func TestAnyOfInlineResponseJSONRoundTrip(t *testing.T) {
	id := "cat-2"
	name := "Luna"
	purrs := true

	var union GetPets200ResponseJSON2
	if err := union.FromCat(Cat{ID: &id, Name: &name, Purrs: &purrs}); err != nil {
		t.Fatalf("FromCat failed: %v", err)
	}

	original := GetPetsJSONResponse{
		Data: []GetPets200ResponseJSON2{union},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded GetPetsJSONResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(decoded.Data) != 1 {
		t.Fatalf("Data length = %d, want 1", len(decoded.Data))
	}

	cat, err := decoded.Data[0].AsCat()
	if err != nil {
		t.Fatalf("AsCat failed: %v", err)
	}
	if *cat.Name != "Luna" {
		t.Errorf("Cat.Name = %q, want %q", *cat.Name, "Luna")
	}
	if *cat.Purrs != true {
		t.Errorf("Cat.Purrs = %v, want true", *cat.Purrs)
	}
}

// TestAnyOfInlineTypeAlias verifies the type alias for the data array.
func TestAnyOfInlineTypeAlias(t *testing.T) {
	var items GetPets200ResponseJSON1
	items = append(items, GetPets200ResponseJSON2{})
	if len(items) != 1 {
		t.Errorf("items length = %d, want 1", len(items))
	}
}

// TestAnyOfInlineApplyDefaults verifies that ApplyDefaults can be called on
// all types without panic.
func TestAnyOfInlineApplyDefaults(t *testing.T) {
	cat := &Cat{}
	cat.ApplyDefaults()

	dog := &Dog{}
	dog.ApplyDefaults()

	rat := &Rat{}
	rat.ApplyDefaults()

	resp := &GetPetsJSONResponse{}
	resp.ApplyDefaults()

	union := &GetPets200ResponseJSON2{}
	union.ApplyDefaults()
}

package card

// Card represents a card entity
// 編號 (ID), 花費 (Cost), 升級費用 (UpgradeCost), 陣營 (Faction), 類別 (Category), 子類別 (SubCategory)
type Card struct {
	ID          int
	Cost        int
	UpgradeCost int
	Faction     string
	Category    string
	SubCategory string
}

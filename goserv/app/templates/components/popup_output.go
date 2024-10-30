package components

type PopupOutput struct {
	JSONButton		Button
	SpreadSheetButton	Button
	ID			string
}

func NewPopupOutput(id string) PopupOutput {
	return PopupOutput {
		JSONButton: Button {
			ID: "JSONButton",
			Text: "JSON format",
			IsSubmit: false,
		},
		SpreadSheetButton: Button {
			ID: "spreadSheetButton",
			Text: "Spreadsheet format",
			IsSubmit: false,
		},
		ID: id,
	}
}

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
			IsPrimary: true,
		},
		SpreadSheetButton: Button {
			ID: "spreadSheetButton",
			Text: "Spreadsheet format",
			IsSubmit: false,
			IsPrimary: true,
		},
		ID: id,
	}
}

package utils


import (
	"fmt"
	"gotransact/apps/transaction/models"
	"gotransact/pkg/db"
	// "log"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func FetchTransactionsLast24Hours() []models.TransactionRequest {
	var transactions []models.TransactionRequest
	last24Hours := time.Now().Add(-24 * time.Hour)
	db.DB.Where("created_at >= ?", last24Hours).Find(&transactions)
	return transactions
}
func GenerateExcel(transactions []models.TransactionRequest) (string, error) {
	f := excelize.NewFile()
	sheetName := "Transactions"
	index := f.NewSheet(sheetName)

	f.SetCellValue(sheetName, "A1", "ID")
	f.SetCellValue(sheetName, "B1", "InternalID")
	f.SetCellValue(sheetName, "C1", "UserID")
	f.SetCellValue(sheetName, "D1", "Status")
	f.SetCellValue(sheetName, "E1", "PaymentGatewayID")
	f.SetCellValue(sheetName, "F1", "Description")
	f.SetCellValue(sheetName, "G1", "Amount")
	f.SetCellValue(sheetName, "H1", "CreatedAt")
	f.SetCellValue(sheetName, "I1", "UpdatedAt")

	for i, tr := range transactions {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), tr.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), tr.InternalID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), tr.UserID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), tr.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), tr.PaymentGatewayMethodID)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), tr.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), tr.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), tr.CreatedAt)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), tr.UpdatedAt)
	}

	f.SetActiveSheet(index)
	filePath := "transaction_report.xlsx"
	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}



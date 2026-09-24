package report

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//go:embed templates/*.html
var templateFS embed.FS

var funcMap = template.FuncMap{
	"now":            time.Now,
	"add":            func(a, b int) int { return a + b },
	"sub":            func(a, b int) int { return a - b },
	"mul":            func(a, b float64) float64 { return a * b },
	"formatMoney":    func(v float64) string { return fmt.Sprintf("%.2f", v) },
	"formatDate":     func(t time.Time) string { return t.Format("02/01/2006") },
	"formatDateTime": func(t time.Time) string { return t.Format("02/01/2006 15:04") },
	"bahtthai":       BahtThai,
}

var templates = map[string]*template.Template{}

func init() {
	names := []string{
		"quotation", "invoice", "purchase_order", "part_picking",
		"receipt", "credit_note", "debit_note", "delivery_sheet",
		"job_card", "daily_sales", "inventory_summary", "customer_list",
	}
	for _, name := range names {
		base := template.New("base.html").Funcs(funcMap)
		tpl := template.Must(base.ParseFS(templateFS, "templates/base.html", "templates/"+name+".html"))
		templates[name] = tpl
	}
}

type TemplateType string

const (
	TplQuotation     TemplateType = "quotation"
	TplInvoice       TemplateType = "invoice"
	TplPurchaseOrder TemplateType = "purchase_order"
	TplPartPicking   TemplateType = "part_picking"
	TplReceipt       TemplateType = "receipt"
	TplCreditNote    TemplateType = "credit_note"
	TplDebitNote     TemplateType = "debit_note"
	TplDeliverySheet TemplateType = "delivery_sheet"
	TplJobCard       TemplateType = "job_card"
	TplDailySales    TemplateType = "daily_sales"
	TplInventorySum  TemplateType = "inventory_summary"
	TplCustomerList  TemplateType = "customer_list"
)

func GenerateHTML(tt TemplateType, data interface{}) (string, error) {
	tmpl, ok := templates[string(tt)]
	if !ok {
		return "", fmt.Errorf("unknown template: %s", tt)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template %s: %w", tt, err)
	}
	return buf.String(), nil
}

func GeneratePDF(ctx context.Context, tt TemplateType, data interface{}) ([]byte, error) {
	html, err := GenerateHTML(tt, data)
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp("", "report-*.html")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := io.WriteString(tmpFile, html); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("write temp file: %w", err)
	}
	tmpFile.Close()

	absPath, _ := filepath.Abs(tmpPath)
	// file:///C:/... on Windows, file:///tmp/... elsewhere; url.URL handles escaping.
	fileURL := (&url.URL{Scheme: "file", Path: "/" + filepath.ToSlash(absPath)}).String()

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
		)...,
	)
	defer allocCancel()

	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	defer tabCancel()

	var pdfBuf []byte
	if err := chromedp.Run(tabCtx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(210.0).
				WithPaperHeight(297.0).
				WithMarginTop(15.0).
				WithMarginBottom(15.0).
				WithMarginLeft(15.0).
				WithMarginRight(15.0).
				Do(ctx)
			return err
		}),
	); err != nil {
		return nil, fmt.Errorf("chromedp pdf: %w", err)
	}

	return pdfBuf, nil
}

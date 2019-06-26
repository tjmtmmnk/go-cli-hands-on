package cmd

import (
	"github.com/spf13/afero"
	"io/ioutil"
	"text/template"
	"time"

	_ "github.com/tjmtmmnk/go-cli-hands-on/statik"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "今日の日報を追加",
	Long:  "テンプレートを元に今日の日報の雛形を作成する",
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName, _ := cmd.Flags().GetString("name")
		_ = generateReport(fileName, afero.NewOsFs())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("name", "n", time.Now().Format("2006-01-02")+".md", "filename")
}

func generateReport(fileName string, afs afero.Fs) error {
	statikFs, _ := fs.New()
	// template読み込む
	tplFile, _ := statikFs.Open("/report.md.tmpl")
	byteTmpl, _ := ioutil.ReadAll(tplFile)
	stringTmpl := string(byteTmpl)
	tmpl := template.Must(template.New("report").Parse(stringTmpl))
	// Todayを差し込む
	// reportFile, _ := os.Create(fileName)
	// afero.Fsは全ての操作を実装していない．そのためAferoに操作を移譲する
	af := afero.Afero{Fs: afs}
	reportFile, _ := af.Create(fileName)
	reportMeta := struct {
		Today string
	}{
		Today: time.Now().Format("2006-01-02"),
	}
	// text/templateとhtml/templateで挙動が違うので注意
	_ = tmpl.Execute(reportFile, reportMeta)
	return nil
}

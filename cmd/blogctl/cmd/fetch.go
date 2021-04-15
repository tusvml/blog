package cmd

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/kjk/notionapi"
	"github.com/km2/notion2hugo"
	"github.com/km2/notion2hugo/frontmatter"
	"github.com/spf13/cobra"
)

var (
	pageID    string
	token     string
	postDir   string
	staticDir string
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch articles from Notion",
	Long:  "fetch articles from Notion.",
	RunE:  runFetch,
}

func runFetch(cmd *cobra.Command, args []string) error {
	client := &notionapi.Client{
		AuthToken: token,
	}

	if err := convertArticles(client, pageID); err != nil {
		return fmt.Errorf("failed to convert: %w", err)
	}

	return nil
}

func convertArticles(c *notionapi.Client, pageID string) error {
	page, err := c.DownloadPage(pageID)
	if err != nil {
		return fmt.Errorf("failed to download page: %w", err)
	}

	if page.TableViews == nil || len(page.TableViews) != 1 {
		return fmt.Errorf("page must have one TableView of articles")
	}

	var result *multierror.Error

	frontMatterConverter := frontmatter.NewConverter(page.TableViews[0]).
		Format(frontmatter.FormatTOML)

	for _, article := range page.TableViews[0].Rows {
		p, err := c.DownloadPage(article.Page.ID)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to download page: %w", err))

			continue
		}

		converter := notion2hugo.NewConverter(c, p).
			FrontMatterConverter(frontMatterConverter).
			PostDir(postDir).
			StaticDir(staticDir)

		if err := converter.ToMarkdownWithFrontMatter(article.Columns, frontmatter.FormatTOML); err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to convert to Markdown: %w", err))

			continue
		}
	}

	if err := result.ErrorOrNil(); err != nil {
		return result
	}

	return nil
}

func init() {
	fetchCmd.Flags().StringVarP(&pageID, "page", "p", "", "Notion page ID with TableView of articles")
	fetchCmd.Flags().StringVarP(&token, "token", "t", "", "Notion token (not required for public pages)")
	fetchCmd.Flags().StringVar(&postDir, "post", "content/posts", "path to the directory that contains articles")
	fetchCmd.Flags().StringVar(&staticDir, "static", "static", "path to the directory that contains static files")

	fetchCmd.MarkFlagRequired("page") //nolint:errcheck

	rootCmd.AddCommand(fetchCmd)
}

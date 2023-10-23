package rtfdoc

import (
	"fmt"
)

// // PageSize returns page Size instance
// func PageSize(width, height int) Size {
// 	return Size{width: width, height: height}
// }

// // PageMargins returns margins
// func PageMargins(left, right, top, bottom int) margins {
// 	return margins{left: left, right: right, top: top, bottom: bottom}
// }

// AddPageHeader returns new PageHeader instance
func (doc *Document) AddPageHeader(headerType PageHeaderType,
	content DocumentItem) *PageHeader {
	pageHead := &PageHeader{
		Type:    headerType,
		Content: content,
	}
	doc.content = append(doc.content, pageHead)
	return pageHead
}

func (ph *PageHeader) compose() string {
	return fmt.Sprintf("{%s%s}", ph.Type, ph.Content.compose())
}

// use \chpgn for page number
// AddPageFooter returns new PageFooter instance
func (doc *Document) AddPageFooter(footerType PageFooterType,
	content DocumentItem) *PageFooter {
	pageFooter := &PageFooter{
		Type:    footerType,
		Content: content,
	}
	doc.content = append(doc.content, pageFooter)
	return pageFooter
}

func (pf *PageFooter) compose() string {
	return fmt.Sprintf("{%s%s}", pf.Type, pf.Content.compose())
}

package service

import (
	"github.com/gofiber/fiber/v2"

	"oapi/conf"
	"oapi/pkg/client/abbr"
)

type Abbr struct {
	cf *conf.Abbr
}

func NewAbbr(cf *conf.Abbr) *Abbr {
	return &Abbr{cf: cf}
}

type Term struct {
	ID                 string `json:"id"`
	Term               string `json:"term"`
	Definition         string `json:"definition"`
	Category           string `json:"category"`
	CategoryName       string `json:"categoryName"`
	ParentCategory     string `json:"parentCategory"`
	ParentCategoryName string `json:"parentCategoryName"`
	Score              string `json:"score"` // float is not good
}

type GetAbbrRes struct {
	Count int     `json:"count"`
	Terms []*Term `json:"terms"`
}

func (svc *Abbr) GetAbbrs(c *fiber.Ctx) error {
	greq := abbr.GetAbbrsReq{
		UID:     svc.cf.UID,
		TokenID: svc.cf.TokenID,
		Term:    c.Params("term"),
	}
	switch typ := c.Query("type", "reverse"); typ {
	case "exact":
		greq.SearchType = abbr.STExact
	case "reverse":
		greq.SearchType = abbr.STReverse
	default:
		return fiber.NewError(fiber.StatusBadRequest, "invalid type: "+typ)
	}
	gres, err := abbr.GetAbbrs(c.Context(), &greq)
	if err != nil {
		return err
	}

	res := GetAbbrRes{
		Count: len(gres.Result),
	}
	for _, t := range gres.Result {
		res.Terms = append(res.Terms, fromAbbrTerm(t))
	}
	return c.JSON(res)
}

func fromAbbrTerm(t *abbr.Term) (res *Term) {
	if t == nil {
		return nil
	}
	return &Term{
		ID:                 t.ID,
		Term:               t.Term,
		Definition:         t.Definition,
		Category:           t.Category,
		CategoryName:       t.CategoryName,
		ParentCategory:     t.ParentCategory,
		ParentCategoryName: t.ParentCategoryName,
		Score:              t.Score,
	}
}

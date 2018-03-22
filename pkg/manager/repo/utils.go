package repo

import (
	"fmt"
	"net/url"
	"strings"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
)

func GetLabelMapFromRequest(req *pb.DescribeReposRequest) (map[string][]string, error) {
	lbs := req.GetLabel().GetValue()
	//lbs = strings.Replace(lbs, ",", "&", -1)
	lbm, err := url.ParseQuery(lbs)
	if err != nil {
		return nil, err
	}

	return lbm, nil
}

func GetSelectorMapFromRequest(req *pb.DescribeReposRequest) (map[string][]string, error) {
	sls := req.GetSelector().GetValue()
	//sls = strings.Replace(sls, ",", "&", -1)
	slm, err := url.ParseQuery(sls)
	if err != nil {
		return nil, err
	}

	return slm, nil
}

func GenerateSelectQuery(query *db.SelectQuery, table, joinTable, keyField, valueField string, labelMap map[string][]string) *db.SelectQuery {
	for k, lbs := range labelMap {
		for i, lb := range lbs {
			lbs[i] = fmt.Sprintf("'%s'", lb)
		}
		alias := fmt.Sprintf("%s_%s", joinTable, k)
		c := fmt.Sprintf("(%s.%s = '%s' and %s.%s in (%s))", alias, keyField, k, alias, valueField, strings.Join(lbs, ","))

		query = query.JoinAs(joinTable, alias, fmt.Sprintf("%s.repo_id = %s.repo_id", alias, table)).Where(c)
	}

	return query
}

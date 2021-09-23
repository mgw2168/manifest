package v1alpha1

//func (c *applicationOperator) ListApps(conditions *params.Conditions, orderBy string, reverse bool, limit, offset int) (*models.PageableResponse, error) {
//	operatorApps, err := c.listOperatorApps()
//	apps, err := c.listApps(conditions)
//	if err != nil {
//		klog.Error(err)
//		return nil, err
//	}
//	apps = filterApps(apps, conditions)
//	operatorApps = filterOperatorApps(operatorApps, conditions)
//
//	if reverse {
//		sort.Sort(sort.Reverse(HelmApplicationList(apps)))
//	} else {
//		sort.Sort(HelmApplicationList(apps))
//	}
//
//	operatorAppsCount := len(operatorApps)
//	totalCount := len(apps)
//	start, end := (&query.Pagination{Limit: limit, Offset: offset}).GetValidPagination(totalCount)
//	items := make([]interface{}, 0)
//
//	if conditions.Match[CategoryId] == "" {
//		if limit > operatorAppsCount && offset < 2 {
//			apps = apps[start+operatorAppsCount : end]
//		} else if offset > 1 {
//			apps = apps[start-operatorAppsCount : end]
//		} else {
//			apps = apps[start:end]
//		}
//		for i := range operatorApps {
//			versions, err := c.getAppVersionByAppLabel(operatorApps[i].Name)
//			if err != nil && !apierrors.IsNotFound(err) {
//				return nil, err
//			}
//			items = append(items, convertOperatorApp(operatorApps[i], versions))
//		}
//		totalCount = totalCount + operatorAppsCount
//		//}
//	} else if conditions.Match[CategoryId] == "radondb" {
//		for i := range operatorApps {
//			versions, err := c.getAppVersionByAppLabel(operatorApps[i].Name)
//			if err != nil && !apierrors.IsNotFound(err) {
//				return nil, err
//			}
//			items = append(items, convertOperatorApp(operatorApps[i], versions))
//		}
//		return &models.PageableResponse{Items: items, TotalCount: totalCount}, nil
//	}
//
//	for i := range apps {
//		versions, err := c.getAppVersionsByAppId(apps[i].GetHelmApplicationId())
//		if err != nil && !apierrors.IsNotFound(err) {
//			return nil, err
//		}
//		ctg, _ := c.ctgLister.Get(apps[i].GetHelmCategoryId())
//		items = append(items, convertApp(apps[i], versions, ctg, 0))
//	}
//	return &models.PageableResponse{Items: items, TotalCount: totalCount}, nil
//}

package main

func main() {
	/**
	  http://139.198.174.152:8000/kapis/openpitrix.io/v1/apps?orderBy=create_time&paging=limit=13,page=1&conditions=status=active,repo_id=repo-operator&reverse=true

	*/

	//repoId := conditions.Match[RepoId]
	//if strings.HasSuffix(repoId, OperatorAppSuffix) {
	//	// todo operator app
	//	apps, err := c.listOperatorApps()
	//	if err != nil {
	//		klog.Error(err)
	//		return nil, err
	//	}
	//	totalCount := len(apps)
	//	items := make([]interface{}, 0, len(apps))
	//
	//	for i := range apps {
	//		items = append(items, map[string]interface{}{
	//			"name": apps[i].Name,
	//			"owner": apps[i].Spec.Description,
	//			"version": apps[i].Status.LatestVersion,
	//		})
	//	}
	//	return &models.PageableResponse{Items: items, TotalCount: totalCount}, nil
	//}
}

//http://139.198.174.152:8000/kapis/openpitrix.io/v1/apps?orderBy=create_time&paging=limit=13,page=1&conditions=status=active,repo_id=repo-helm&reverse=true

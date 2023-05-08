package navigator

// func TestNavigator_Next_SegmentNotFound(t *testing.T) {
// 	nav, err := New()
// 	require.NoError(t, err)
// 	err = nav.AddRouteFromPoints([][]Point{
// 		{
// 			{X: 62.0084382486728, Y: 108.35674751317089},
// 			{X: 62.004805686033905, Y: 108.36984468615873},
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	nav.Next(1, nav.TotalDistance()*2)
// 	require.Zero(t, nav.RouteIndex())
// 	require.Zero(t, nav.SegmentIndex())
// 	require.Zero(t, nav.TrackIndex())
// }

// func TestNavigator_Next_SegmentFound(t *testing.T) {
// 	r := loadData("short_path")
// 	routes, err := RoutesFromGeoJSON(r)
// 	require.NoError(t, err)
// 	nav, err := New()
// 	require.NoError(t, err)
// 	nav.AddRoutes(routes)

// 	nav.segmentDistance = 800
// 	nav.Next(1, 100)
// 	nav.segmentDistance = 900
// 	nav.Next(1, 200)
// 	require.Equal(t, 0, nav.RouteIndex())
// 	require.Equal(t, 1, nav.SegmentIndex())
// 	require.Equal(t, 0, nav.TrackIndex())
// }

// func TestNavigator_Next_SegmentLoop(t *testing.T) {
// 	r := loadData("short_path")
// 	routes, err := RoutesFromGeoJSON(r)
// 	require.NoError(t, err)
// 	nav, err := New()
// 	require.NoError(t, err)
// 	nav.AddRoutes(routes)
// 	require.Len(t, routes, 1)

// 	for i := 0; i < 8; i++ {
// 		nav.segmentDistance = 900
// 		nav.Next(1, 10)
// 		require.Equal(t, 1, nav.SegmentIndex())
// 		require.Equal(t, 0, nav.RouteIndex())
// 	}
// }

// func TestNavigator_Next_NextRoute(t *testing.T) {
// 	r := loadData("routes_2_linestring")
// 	routes, err := RoutesFromGeoJSON(r)
// 	require.NoError(t, err)
// 	require.Len(t, routes, 2)

// 	nav, err := New()
// 	require.NoError(t, err)
// 	nav.AddRoutes(routes)
// 	fmt.Println(nav.TotalDistance())

// 	for i := 0; i < int(nav.TotalDistance())*2; i += 100 {
// 		ok := nav.Next(1, 100)
// 		fmt.Println("next", ok)
// 		fmt.Println(nav.RouteIndex(), nav.TrackIndex(), nav.SegmentIndex(), nav.CurrentDistance())
// 		fmt.Println(nav.Location())
// 	}
// }

package rest

//func authMiddleware(c *gin.Context) {
//
//	userID := c.Header.Get("id-x")
//
//	if userID == "" {
//		log.Printf("[%s], %s - error: userID not provided", r.Method, r.URL)
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	ctx := r.Context()
//	ctx = context.WithValue(ctx, "id", userID)
//	r = r.WithContext(ctx)
//	c.Next()
//}
//
//func loggerMiddleware(c *gin.Context) {
//	idFromCtx := r.Context().Value("id")
//	id, ok := idFromCtx.(string)
//	if !ok {
//		log.Printf("[%s], %s - error: userID is invalid", r.Method, r.URL)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	log.Printf("[%s], %s by userId %s", r.Method, r.URL, id)
//	c.Next()
//}

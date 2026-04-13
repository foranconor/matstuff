package main

import (
	"github.com/joho/godotenv"
	"log"
	"matstuff/db"
	"matstuff/handlers"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set in .env")
	}

	database, err := db.Open()
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	auth, err := handlers.NewAuth()
	if err != nil {
		log.Fatal("failed to initialise auth:", err)
	}

	h := &handlers.Handler{DB: database}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/nonce", auth.NonceRoute)
	mux.HandleFunc("POST /api/auth/tokens", auth.TokensRoute)

	// Materials
	mux.HandleFunc("POST /api/materials", h.CreateMaterial)
	mux.HandleFunc("GET /api/materials", h.ListMaterials)
	mux.HandleFunc("GET /api/materials/{id}", h.GetMaterial)
	mux.HandleFunc("PUT /api/materials/{id}", h.UpdateMaterial)
	mux.HandleFunc("PUT /api/materials/{id}/uses", h.UpdateUses)
	mux.HandleFunc("GET /api/materials/{id}/suppliers", h.ListMaterialSuppliers)
	mux.HandleFunc("POST /api/materials/{id}/suppliers", h.AddSupplierToMaterial)
	mux.HandleFunc("GET /api/materials/{id}/profiles", h.ListMaterialProfiles)
	mux.HandleFunc("POST /api/materials/{id}/profiles", h.AddProfileToMaterial)
	mux.HandleFunc("GET /api/materials/{id}/notes", h.ListMaterialNotes)
	mux.HandleFunc("POST /api/materials/{id}/notes", h.CreateMaterialNote)
	mux.HandleFunc("PUT /api/materials/{id}/notes/{nid}", h.UpdateMaterialNote)
	mux.HandleFunc("DELETE /api/materials/{id}/notes/{nid}", h.DeleteMaterialNote)

	// Suppliers
	mux.HandleFunc("POST /api/suppliers", h.CreateSupplier)
	mux.HandleFunc("GET /api/suppliers", h.ListSuppliers)
	mux.HandleFunc("GET /api/suppliers/{id}", h.GetSupplier)
	mux.HandleFunc("GET /api/suppliers/{id}/materials", h.ListSupplierMaterials)
	mux.HandleFunc("POST /api/suppliers/{id}/materials", h.AddMaterialToSupplier)
	mux.HandleFunc("PUT /api/suppliers/{id}", h.UpdateSupplier)
	mux.HandleFunc("PUT /api/supplier-materials/{id}", h.UpdateMaterialSupplier)
	mux.HandleFunc("DELETE /api/supplier-materials/{id}", h.DeleteMaterialSupplier)

	// Profiles
	mux.HandleFunc("GET /api/profiles", h.ListProfiles)
	mux.HandleFunc("POST /api/profiles", h.CreateProfile)
	mux.HandleFunc("GET /api/profiles/{id}", h.GetProfile)
	mux.HandleFunc("PUT /api/profiles/{id}", h.UpdateProfile)
	mux.HandleFunc("DELETE /api/profiles/{id}", h.DeleteProfile)
	mux.HandleFunc("GET /api/profiles/{id}/materials", h.ListProfileMaterials)
	mux.HandleFunc("POST /api/profiles/{id}/materials", h.AddMaterialToProfile)
	mux.HandleFunc("PUT /api/material-profiles/{id}", h.UpdateMaterialProfile)
	mux.HandleFunc("DELETE /api/material-profiles/{id}", h.DeleteMaterialProfile)

	// Contacts
	mux.HandleFunc("POST /api/contacts", h.CreateContact)
	mux.HandleFunc("GET /api/contacts", h.ListContacts)
	mux.HandleFunc("GET /api/contacts/{id}", h.GetContact)
	mux.HandleFunc("PUT /api/contacts/{id}", h.UpdateContact)

	// Brackets
	mux.HandleFunc("POST /api/brackets", h.CreateBracket)
	mux.HandleFunc("GET /api/brackets", h.ListBrackets)
	mux.HandleFunc("GET /api/brackets/{id}", h.GetBracket)
	mux.HandleFunc("PUT /api/brackets/{id}", h.UpdateBracket)
	mux.HandleFunc("GET /api/brackets/{id}/notes", h.ListBracketNotes)
	mux.HandleFunc("POST /api/brackets/{id}/notes", h.CreateBracketNote)
	mux.HandleFunc("PUT /api/brackets/{id}/notes/{nid}", h.UpdateBracketNote)
	mux.HandleFunc("DELETE /api/brackets/{id}/notes/{nid}", h.DeleteBracketNote)

	// Serve static frontend files, falling back to index.html for SPA routing
	fs := http.FileServer(http.Dir("web/dist"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "web/dist" + r.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, "web/dist/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	handler := Middleware(mux,
		AuthMiddleware(auth),
		CorsMiddleware,
		RouteLoggerMiddleware,
		RecoveryMiddleware,
	)

	log.Println("Listening on :" + port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

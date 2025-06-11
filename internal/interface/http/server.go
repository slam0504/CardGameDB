package http

import (
	stdhttp "net/http"
	"strconv"

	"CardGameDB/internal/domain/card"
	"CardGameDB/internal/infrastructure/eventbus"
)

const (
	searchEvent = "card.search"
	createEvent = "card.create"
	updateEvent = "card.update"
)

// Server struct

type Server struct {
	bus *eventbus.EventBus
}

// NewServer creates server
func NewServer(bus *eventbus.EventBus) *Server {
	return &Server{bus: bus}
}

// Start HTTP server on addr
func (s *Server) Start(addr string) error {
	stdhttp.HandleFunc("/search", s.handleSearch)
	stdhttp.HandleFunc("/create", s.handleCreate)
	stdhttp.HandleFunc("/update", s.handleUpdate)
	return stdhttp.ListenAndServe(addr, nil)
}

func (s *Server) handleSearch(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var filter card.Filter

	if v := r.URL.Query().Get("id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filter.ID = &id
		}
	}
	if v := r.URL.Query().Get("cost"); v != "" {
		if c, err := strconv.Atoi(v); err == nil {
			filter.Cost = &c
		}
	}
	if v := r.URL.Query().Get("upgrade_cost"); v != "" {
		if uc, err := strconv.Atoi(v); err == nil {
			filter.UpgradeCost = &uc
		}
	}
	if v := r.URL.Query().Get("faction"); v != "" {
		filter.Faction = &v
	}
	if v := r.URL.Query().Get("category"); v != "" {
		filter.Category = &v
	}
	if v := r.URL.Query().Get("subcategory"); v != "" {
		filter.SubCategory = &v
	}

	reply := make(chan card.SearchResult)
	s.bus.Publish(searchEvent, card.SearchRequested{Filter: filter, Reply: reply})
	res := <-reply
	if res.Err != nil {
		stdhttp.Error(w, res.Err.Error(), stdhttp.StatusInternalServerError)
		return
	}

	// simple JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("["))
	for i, c := range res.Cards {
		if i > 0 {
			w.Write([]byte(","))
		}
		w.Write([]byte("{"))
		w.Write([]byte("\"id\":" + strconv.Itoa(c.ID)))
		w.Write([]byte(",\"cost\":" + strconv.Itoa(c.Cost)))
		w.Write([]byte(",\"upgrade_cost\":" + strconv.Itoa(c.UpgradeCost)))
		w.Write([]byte(",\"faction\":\"" + c.Faction + "\""))
		w.Write([]byte(",\"category\":\"" + c.Category + "\""))
		w.Write([]byte(",\"subcategory\":\"" + c.SubCategory + "\""))
		w.Write([]byte("}"))
	}
	w.Write([]byte("]"))
}

func (s *Server) handleCreate(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodPost {
		stdhttp.Error(w, "method not allowed", stdhttp.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	cost, _ := strconv.Atoi(r.FormValue("cost"))
	uc, _ := strconv.Atoi(r.FormValue("upgrade_cost"))
	cardData := card.Card{
		ID:          id,
		Cost:        cost,
		UpgradeCost: uc,
		Faction:     r.FormValue("faction"),
		Category:    r.FormValue("category"),
		SubCategory: r.FormValue("subcategory"),
	}

	reply := make(chan error)
	s.bus.Publish(createEvent, card.CreateRequested{Card: cardData, Reply: reply})
	err := <-reply
	if err != nil {
		stdhttp.Error(w, err.Error(), stdhttp.StatusInternalServerError)
		return
	}
	w.WriteHeader(stdhttp.StatusCreated)
}

func (s *Server) handleUpdate(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	if r.Method != stdhttp.MethodPost && r.Method != stdhttp.MethodPut {
		stdhttp.Error(w, "method not allowed", stdhttp.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))
	cost, _ := strconv.Atoi(r.FormValue("cost"))
	uc, _ := strconv.Atoi(r.FormValue("upgrade_cost"))
	cardData := card.Card{
		ID:          id,
		Cost:        cost,
		UpgradeCost: uc,
		Faction:     r.FormValue("faction"),
		Category:    r.FormValue("category"),
		SubCategory: r.FormValue("subcategory"),
	}

	reply := make(chan error)
	s.bus.Publish(updateEvent, card.UpdateRequested{Card: cardData, Reply: reply})
	err := <-reply
	if err != nil {
		stdhttp.Error(w, err.Error(), stdhttp.StatusInternalServerError)
		return
	}
	w.WriteHeader(stdhttp.StatusOK)
}

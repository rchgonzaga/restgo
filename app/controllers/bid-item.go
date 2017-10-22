package controllers

import (
	"encoding/json"
	"restgo/app/models"

	"github.com/revel/revel"
)

type BidItemCtrl struct {
	GorpController
}

func (bidItemCtrl BidItemCtrl) parseBidItem() (models.BidItem, error) {
	biditem := models.BidItem{}
	err := json.NewDecoder(bidItemCtrl.Request.Body).Decode(&biditem)
	return biditem, err
}

// Add a new item bid
func (bidItemCtrl BidItemCtrl) Add() revel.Result {
	if biditem, err := bidItemCtrl.parseBidItem(); err != nil {
		return bidItemCtrl.RenderText("Unable to parse the BidItem from JSON.")
	} else {
		// Validate the model
		biditem.Validate(bidItemCtrl.Validation)
		if bidItemCtrl.Validation.HasErrors() {
			// Do something better here!
			return bidItemCtrl.RenderText("You have error in your BidItem.")
		} else {
			if err := bidItemCtrl.Txn.Insert(&biditem); err != nil {
				return bidItemCtrl.RenderText(
					"Error inserting record into database!")
			} else {
				return bidItemCtrl.RenderJSON(biditem)
			}
		}
	}
}

// Get a bid item
func (bidItemCtrl BidItemCtrl) Get(id int64) revel.Result {
	biditem := new(models.BidItem)
	err := bidItemCtrl.Txn.SelectOne(biditem,
		`SELECT * FROM BidItem WHERE id = ?`, id)
	if err != nil {
		return bidItemCtrl.RenderText("Error.  Item probably doesn't exist.")
	}
	return bidItemCtrl.RenderJSON(biditem)
}

// List all items in bid table
func (bidItemCtrl BidItemCtrl) List() revel.Result {
	lastID := parseIntOrDefault(bidItemCtrl.Params.Get("lid"), -1)
	limit := parseUintOrDefault(bidItemCtrl.Params.Get("limit"), uint64(25))
	biditems, err := bidItemCtrl.Txn.Select(models.BidItem{},
		`SELECT * FROM BidItem WHERE Id > ? LIMIT ?`, lastID, limit)
	if err != nil {
		return bidItemCtrl.RenderText(
			"Error trying to get records from DB.")
	}
	return bidItemCtrl.RenderJSON(biditems)
}

// Update a specific item
func (bidItemCtrl BidItemCtrl) Update(id int64) revel.Result {
	biditem, err := bidItemCtrl.parseBidItem()
	if err != nil {
		return bidItemCtrl.RenderText("Unable to parse the BidItem from JSON.")
	}
	// Ensure the Id is set.
	biditem.Id = id
	success, err := bidItemCtrl.Txn.Update(&biditem)
	if err != nil || success == 0 {
		return bidItemCtrl.RenderText("Unable to update bid item.")
	}
	return bidItemCtrl.RenderText("Updated %v", id)
}

// Delete a specific item
func (bidItemCtrl BidItemCtrl) Delete(id int64) revel.Result {
	success, err := bidItemCtrl.Txn.Delete(&models.BidItem{Id: id})
	if err != nil || success == 0 {
		return bidItemCtrl.RenderText("Failed to remove BidItem")
	}
	return bidItemCtrl.RenderText("Deleted %v", id)
}

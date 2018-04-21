(ns bitkit.views
  (:require [re-frame.core :as re-frame]
            [reagent.core :as reagent]
            [bitkit.subs :as subs]
            [bitkit.events :as events]
            [bitkit.routes :refer [set-path!]]))

(defn handler
  "Wraps an event handler function so that it first
  prevents the event default"
  [func]
  (fn [event]
    (.preventDefault event)
    (func event)))

(defn transaction-form
  [{:keys [txid]}]
  [:form.content
   {:on-submit
    (handler #(re-frame/dispatch [::events/set-transaction txid]))}
   [:div.field
    [:label.label "Transaction ID"]
    [:div.control
     [:input.input
      {:value     txid
       :on-change #(re-frame/dispatch [::events/set-transaction-id (.. % -target -value)])}]]]])

(defn notification
  [{:keys [error]}]
  (when error
    [:div.notification.is-warning.content
     [:p
      "The given transaction ID could not be found in the mempool.
      This can happen for a variety of reasons:"]
     [:ul
      [:li "The transaction ID was entered incorrectly"]
      [:li "The transaction has already been confirmed"]
      [:li "The transaction has been evicted"]]]))

(defn transaction
  [{:keys [txn]}]
  (when txn
    [:section
     [:div.content.is-small
      [:h2 "Your transaction"]
      [:ul.is-marginless.is-unstyled
       [:li (str "Fee: " (:fee txn) " satoshis")]
       [:li (str "Fee rate: " (:fee_rate txn) " satoshis per byte")]]]
     [:div.content.is-small
      [:h2 "Transactions with a higher fee rate"]
      [:ul.is-marginless.is-unstyled
       [:li (str "Count: " (:transaction_count txn))]
       [:li (str "Block capacity used: " (:capacity_used txn))]]]]))

(defn main-panel []
  (let [txid  (re-frame/subscribe [::subs/transaction-id])
        txn   (re-frame/subscribe [::subs/transaction])
        error (re-frame/subscribe [::subs/error])]
    [:section.section
     [:div.container
      [transaction-form {:txid @txid}]
      [notification {:error @error}]
      [transaction {:txn @txn}]]]))

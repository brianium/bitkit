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
  [:form
   {:on-submit
    (handler #(re-frame/dispatch [::events/set-transaction txid]))}
   [:div.field
    [:div.transaction-label
     [:label.label.is-medium "Transaction ID"]
     [:button.button.is-info
      {:type     "button"
       :on-click #(re-frame/dispatch [::events/random-transaction])}
      "random"]]
    [:div.control
     [:input.input.is-large
      {:value     txid
       :on-change #(re-frame/dispatch [::events/set-transaction-id (.. % -target -value)])}]]]])

(defn notification
  [{:keys [error]}]
  (when error
    (case error
      :left-mempool
      [:div.notification.is-success.content
       [:p "Your transaction has left the mempool!"]]

      :not-found
      [:div.notification.is-warning.content
       [:p
        "The given transaction ID could not be found in the mempool.
      This can happen for a variety of reasons:"]
       [:ul
        [:li "The transaction ID was entered incorrectly"]
        [:li "The transaction has already been confirmed"]
        [:li "The transaction has been evicted"]]]

      [:div.notification.is-danger.content
       :p "An unknown error occurred. Sorry!"])))

(defn interval-list
  [& children]
  (let [ref (atom nil)]
    (reagent/create-class
      {:component-did-update
       (fn []
         (-> @ref
             .-classList
             (.add "is-updated"))
         (js/setTimeout
           (fn []
             (-> @ref
                 .-classList
                 (.remove "is-updated")))
           1500))
       :reagent-render
       (fn [& children]
         [:div.interval-list {:ref (fn [elem] (reset! ref elem))}
          [:ul
           {:class-name (str "is-marginless is-unstyled is-size-6")}
           (map-indexed #(with-meta %2 {:key %1}) children)]])})))

(defn transaction
  [{:keys [txn]}]
  (when txn
    [:section
     [:div.message.is-info
      [:h2.message-header "Your transaction"]
      [interval-list
       [:div.message-body
        [:li (str "Fee: " (:fee txn) " satoshis")]
        [:li (str "Fee rate: " (:fee_rate txn) " satoshis per vbyte")]
        [:li (str "Virtual size: " (:weight txn) " vbytes")]]]]
     [:div.message.is-info
      [:h2.message-header "Transactions with a higher fee rate"]
      [interval-list
       [:div.message-body
        [:li (str "Count: " (:transaction_count txn))]
        [:li (str "Block capacity used: " (:capacity_used txn))]]]]
     [:div.message.is-info
      [:h2.message-header "Entire mempool"]
      [interval-list
       [:div.message-body
        [:li (str "Count: " (:mempool_transaction_count txn))]
        [:li (str "Virtual size: " (:mempool_total_virtual_size txn) " vbytes")]]]]]))

(defn main-panel []
  (let [txid  (re-frame/subscribe [::subs/transaction-id])
        txn   (re-frame/subscribe [::subs/transaction])
        error (re-frame/subscribe [::subs/error])]
    [:section.section.main-panel
     [:div.container
      [transaction-form {:txid @txid}]
      [notification {:error @error}]
      [transaction {:txn @txn}]]]))

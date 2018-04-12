(ns bitkit.views
  (:require [re-frame.core :as re-frame]
            [reagent.core :as reagent]
            [bitkit.subs :as subs]
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
  (let [value (reagent/atom txid)]
    (fn [{:keys [txid]}]
      [:form {:on-submit (handler #(set-path! (str "/" @value)))}
       [:div.field
        [:label.label "Transaction ID"]
        [:div.control
         [:input.input
          {:value     @value
           :on-change #(reset! value (.. % -target -value))}]]
        [:p.help "Bitcoin transaction id"]]])))

(defn notification
  [{:keys [error]}]
  (when error
    [:div.notification.is-warning
     "The given transaction ID could not be found in the mempool. This means
     it has already been confirmed or never existed in the first place. We should
     totally work on the wording/experience of this message"]))

(defn transaction
  [{:keys [txn]}]
  (when txn
    [:div.has-text-centered
     [:p
      [:span.is-size-3 (:transaction_count txn)]
      " transactions ahead of you"]
     [:p (str "totaling " (:total_weight txn) " whatever units ryan tell me")]]))

(defn main-panel []
  (let [txid (re-frame/subscribe [::subs/transaction-id])
        txn  (re-frame/subscribe [::subs/transaction])
        error (re-frame/subscribe [::subs/error])]
    [:section.section
     [:div.container
      [transaction-form {:txid @txid}]
      [notification {:error @error}]
      [transaction {:txn @txn}]]]))

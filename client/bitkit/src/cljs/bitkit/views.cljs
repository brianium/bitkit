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

(defn transaction-form []
  (let [txid  @(re-frame/subscribe [::subs/transaction-id])
        value (reagent/atom txid)]
    (fn [props]
      [:form {:on-submit (handler #(set-path! (str "/" @value)))}
       [:div.field
        [:label.label "Transaction ID"]
        [:div.control
         [:input.input
          {:value     @value
           :on-change #(reset! value (.. % -target -value))}]]
        [:p.help "Bitcoin transaction id"]]])))

(defn main-panel []
  [:section.section
   [:div.container
    [transaction-form]]])

(ns bitkit.events
  (:require [re-frame.core :as re-frame]
            [bitkit.db :as db]
            [ajax.core :as ajax]))

(re-frame/reg-event-db
 ::initialize-db
 (fn  [_ _]
   db/default-db))

(defn transaction
  "Takes a transaction id and updates state with transaction
  data"
  [cofx id]
  {:db         (assoc-in cofx [:db :transaction-id] id)
   :http-xhrio {:method          :get
                :uri             (str "https://api.bitkit.live/transactions/" id)
                :response-format (ajax/json-response-format {:keywords? true})
                :on-success      [:fetch-transaction-success]
                :on-failure      [:fetch-transaction-error]}})

(defn index
  [cofx]
  {:db (assoc cofx :db db/default-db)})

(re-frame/reg-event-fx
  ::set-route
  (fn [cofx [_ {:keys [route-params handler]}]]
    (case handler
      :transaction (transaction cofx (:id route-params))
      (index cofx))))

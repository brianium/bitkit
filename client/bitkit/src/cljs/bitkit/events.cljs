(ns bitkit.events
  (:require [re-frame.core :as re-frame]
            [bitkit.db :as db]
            [bitkit.routes :as routes]
            [ajax.core :as ajax]))

(re-frame/reg-event-db
 ::initialize-db
 (fn  [_ _]
   db/default-db))

(re-frame/reg-event-fx
  ::set-transaction-id
  (fn [{:keys [db]} [_ id]]
    {:db                    (assoc db :transaction-id id)
     ::transaction-interval {:action      :stop
                             :interval-id (:interval db)}}))

(re-frame/reg-event-db
  ::set-interval
  (fn [db [_ id]]
    (assoc db :interval id)))

(defn transaction
  "Takes a transaction id and updates state with transaction
  data"
  [{:keys [db]} id]
  {:http-xhrio {:method          :get
                :uri             (str "https://api.bitkit.live/transactions/" id)
                :response-format (ajax/json-response-format {:keywords? true})
                :on-success      [::fetch-transaction-success]
                :on-failure      [::fetch-transaction-error (:interval db)]}
   :dispatch   [::set-transaction-id id]})

(defn index
  [cofx]
  {:db (assoc cofx :db db/default-db)})

(re-frame/reg-event-fx
  ::set-route
  (fn [cofx [_ {:keys [route-params handler]}]]
    (case handler
      :transaction (transaction cofx (:id route-params))
      (index cofx))))

(re-frame/reg-event-fx
  ::update-transaction
  (fn [cofx [_ txid]]
    (transaction cofx txid)))

(re-frame/reg-event-fx
  ::set-transaction
  (fn [cofx [_ txid]]
    {:db            (:db cofx)
     ::set-tx-route txid}))

(re-frame/reg-event-fx
  ::fetch-transaction-success
  (fn [{:keys [db]} [_ response]]
    {:db                    (-> db
                                (assoc :transaction (:data response))
                                (assoc :error nil))
     ::transaction-interval {:previous-txid (:transaction-id db)
                             :txid          (get-in response [:data :txid])
                             :action        :start
                             :interval-id   (:interval db)}}))

(re-frame/reg-event-fx
  ::fetch-transaction-error
  (fn [_ [_ interval-id]]
    {:db
     (merge db/default-db {:error
                           (if interval-id :left-mempool :not-found)})
     
     ::transaction-interval
     {:action      :stop
      :interval-id interval-id}}))

;;; Side effects

(defn set-tx-effect
  "Handles updating the current transaction in scope. Takes a
  transaction id"
  [txid]
  (routes/set-path! (str "/" txid)))

(re-frame/reg-fx ::set-tx-route set-tx-effect)

(defn stop-interval
  [{:keys [interval-id]}]
  (when interval-id
    (js/clearInterval interval-id)
    (re-frame/dispatch [::set-interval nil])))

(defn start-interval
  [{:keys [txid previous-txid interval-id]}]
  (-> #(re-frame/dispatch [::update-transaction txid])
      (js/setInterval 5000)
      (as-> id [::set-interval id])
      re-frame/dispatch))

(def restart-interval (juxt stop-interval start-interval))

(defn update-interval
  [{:keys [txid previous-txid action interval-id] :as params}]
  (case action
    :start (if (or (nil? interval-id) (not= txid previous-txid))
             (start-interval params)
             (restart-interval params))
    (stop-interval params)))

(re-frame/reg-fx ::transaction-interval update-interval)

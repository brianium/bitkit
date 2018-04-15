(ns bitkit.subs
  (:require [re-frame.core :as re-frame]))

(def blocksize-in-bytes 4194304)

(re-frame/reg-sub
 ::transaction-id
 (fn [db]
   (or (:transaction-id db) "")))

(re-frame/reg-sub
  ::error
  (fn [db]
    (:error db)))

(re-frame/reg-sub
  ::transaction
  (fn [{:keys [transaction]}]
    (when transaction
      (let [{:keys [fee_rate, weight, total_weight]} transaction]
        (-> transaction
            (assoc :fee (* fee_rate weight))
            (assoc :capacity_used (/ total_weight blocksize-in-bytes)))))))

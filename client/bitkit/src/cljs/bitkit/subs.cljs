(ns bitkit.subs
  (:require [re-frame.core :as re-frame]
            [goog.string :as gstring]
            [goog.string.format]))

(def blocksize-in-bytes 1000000)

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
      (let [{:keys [fee, weight, total_weight]} transaction]
        (-> transaction
            (assoc :fee_rate (gstring/format "%.1f" (/ fee weight)))
            (assoc :capacity_used
              (-> (/ total_weight blocksize-in-bytes)
                  (* 100)
                  (as-> capacity (gstring/format "%.1f" capacity))
                  (str "%"))))))))

(ns bitkit.subs
  (:require [re-frame.core :as re-frame]))

(re-frame/reg-sub
 ::transaction-id
 (fn [db]
   (or (:transaction-id db) "")))

(re-frame/reg-sub
  ::transaction
  (fn [db]
    (:transaction db)))

(re-frame/reg-sub
  ::error
  (fn [db]
    (:error db)))

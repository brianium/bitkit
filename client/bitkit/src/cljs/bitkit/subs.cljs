(ns bitkit.subs
  (:require [re-frame.core :as re-frame]))

(re-frame/reg-sub
 ::transaction-id
 (fn [db]
   (or (:transaction-id db) "")))

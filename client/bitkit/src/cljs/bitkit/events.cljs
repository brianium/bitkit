(ns bitkit.events
  (:require [re-frame.core :as re-frame]
            [bitkit.db :as db]))

(re-frame/reg-event-db
 ::initialize-db
 (fn  [_ _]
   db/default-db))

(re-frame/reg-event-db
  :set-route
  (fn [db [_ match]]
    (.log js/console (clj->js match))))

(ns bitkit.core
  (:require [reagent.core :as reagent]
            [re-frame.core :as re-frame]
            [day8.re-frame.http-fx]
            [bitkit.events :as events]
            [bitkit.views :as views]
            [bitkit.config :as config]
            [bitkit.routes :as routes]))

(defn dev-setup []
  (when config/debug?
    (enable-console-print!)
    (println "dev mode")))

(defn mount-root []
  (re-frame/clear-subscription-cache!)
  (reagent/render [views/main-panel]
                  (.getElementById js/document "app")))

(defn ^:export init []
  (re-frame/dispatch-sync [::events/initialize-db])
  (routes/listen! ::events/set-route)
  (dev-setup)
  (mount-root))

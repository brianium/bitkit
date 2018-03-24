(ns bitkit.views
  (:require [re-frame.core :as re-frame]
            [bitkit.subs :as subs]
            ))

(defn main-panel []
  (let [name (re-frame/subscribe [::subs/name])]
    [:div "Hello there from " @name]))

# Gehirnstürmen

## Eingeklappt

- Kontroll-Knöpfe (Start, Stop, Neustart, Stop)
- Stapel und Behälter Status (klein)
- Stapel Name
- Netzbenutzerimgesicht Knöpfe (wenn Etikett gesetzt)
- Link zu Portainer zum bearbeiten


## Ausgeklappt

- Verbundene Netzwerke
- Gekartete Anschlüsse
- Einzelinformationen über die Behälter
-



## Konzept?

Wir brauchen eine Gruppenkomponente, und diese Komponente ist ein Stapel-Behälter, sie verhält sich so ähnlich wie ein Stapel Behälter, aber innerhalb davon ist es wieder eine sortierte Stapel Tabelle.
Eine Gruppe, wenn ausgeklappt, zeigt die Stapel von dieser Gruppe an und die Stapel sind eingerückt. Gruppe ist in einer Karte/Bettlaken drin, sodass man sieht, dass es einfach in dem einem Behälter drin ist



Wir gebe eine Liste von Stapeln hinein, wenn es mehr als ein Stapel ist dann behandeln wir es als Gruppe!




```js

class SortableMasterUnit {

    constructor(stacks: Stack[])
}




[
    {
        type: "group",
        group: {
            groupName: string;
            globalPriority: number;
            stacks: [stack1, stack2]
        }
        stack: undefined
    },
    {
        type: "group",
        group: {
            groupName: string;
            globalPriority: number;
            stacks: [stack1, stack2]
        }
        stack: undefined
    },
    {
        type: "group",
        group: {
            groupName: string;
            globalPriority: number;
            stacks: [stack1, stack2]
        }
        stack: undefined
    },
    {
        type: "stack"
        group: undefined
        stack: Stack
    }
]
```

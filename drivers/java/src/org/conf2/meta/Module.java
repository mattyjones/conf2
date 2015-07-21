package org.conf2.meta;

/**
 *
 */
public class Module extends CollectionBase implements HasChoices, HasGroupings, HasTypedefs, Describable {
    private String namespace;
    private String prefix;
    private LinkedListCollection groupings = new LinkedListCollection("groupings", this);
    private LinkedListCollection choices = new LinkedListCollection("choices", this);
    private LinkedListCollection typedefs = new LinkedListCollection("typedefs", this);
    private LinkedListCollection notifications = new LinkedListCollection("notifications", this);
    private LinkedListCollection rpcs = new LinkedListCollection("rpcs", this);

    public Module(String ident) {
        super(ident);
    }

    public void setNamespace(String namespace) {
        this.namespace = namespace;
    }

    public String getNamespace() {
        return namespace;
    }

    public void setPrefix(String prefix) {
        this.prefix = prefix;
    }

    public String getPrefix() {
        return prefix;
    }

    @Override
    public void addMeta(Meta m) {
        Class c = m.getClass();
        if (c == Grouping.class) {
            groupings.addMeta(m);
        } else {
            super.addMeta(m);
        }
    }

    public MetaCollection getRpcs() {
        return rpcs;
    }

    public MetaCollection getNotifications() {
        return notifications;
    }

    public MetaCollection getChoices() {
        return choices;
    }

    public MetaCollection getTypedefs() {
        return typedefs;
    }

    public MetaCollection getGroupings() {
        return groupings;
    }
}

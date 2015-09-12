package org.conf2.schema;

/**
 *
 */
public class Container extends CollectionBase implements Describable, HasTypedefs, HasGroupings {
    private LinkedListCollection groupings = new LinkedListCollection("groupings", this);
    private LinkedListCollection typedefs = new LinkedListCollection("typedefs", this);
    private boolean config;
    private boolean mandatory;

    public Container(String ident) {
        super(ident);
    }

    @Override
    public void addMeta(Meta m) {
        if (m instanceof Grouping) {
            groupings.addMeta(m);
        } else if (m instanceof Typedef) {
            typedefs.addMeta(m);
        } else {
            super.addMeta(m);
        }
    }

    @Override
    public MetaCollection getGroupings() {
        return groupings;
    }

    @Override
    public MetaCollection getTypedefs() {
        return typedefs;
    }

    public void setConfig(boolean config) {
        this.config = config;
    }

    public void setMandatory(boolean mandatory) {
        this.mandatory = mandatory;
    }
}

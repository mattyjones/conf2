package org.conf2.schema;

/**
 *
 */
public class MetaList extends CollectionBase implements Describable, HasTypedefs, HasGroupings {
    private LinkedListCollection groupings = new LinkedListCollection("groupings", this);
    private LinkedListCollection typedefs = new LinkedListCollection("typedefs", this);
    private String[] keys;
    private boolean config;
    private boolean mandatory;

    public MetaList(String ident) {
        super(ident);
    }

    public void setEncodedKeys(String encodedKeys) {
        keys = encodedKeys.split(" ");
    }

    public String[] getKeys() {
        return keys;
    }

    public void setKeys(String[] keys) {
        this.keys = keys;
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

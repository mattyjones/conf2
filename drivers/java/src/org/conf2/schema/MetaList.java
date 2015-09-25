package org.conf2.schema;

/**
 *
 */
public class MetaList extends CollectionBase implements Describable, HasTypedefs, HasGroupings {
    private LinkedListCollection groupings = new LinkedListCollection("groupings", this);
    private LinkedListCollection typedefs = new LinkedListCollection("typedefs", this);
    private String[] keys;
    private DataType[] keyTypes;
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

    public DataType[] getKeyDataTypes() {
        if (keyTypes == null) {
            DataType[] types = new DataType[this.keys.length];
            for (int i = 0; i < types.length; i++) {
                HasDataType m = (HasDataType) MetaUtil.findByIdent(this, this.keys[i]);
                types[i] = m.getDataType();
            }
            keyTypes = types;
        }
        return keyTypes;
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

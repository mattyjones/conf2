package org.conf2.yang;

/**
 *
 */
public class Typedef extends MetaBase implements Describable {
    private DataType dataType;

    public Typedef(String ident) {
        super(ident);
    }

    public void setDataType(DataType dataType) {
        this.dataType = dataType;
    }

    public DataType getDataType() {
        return dataType;
    }
}

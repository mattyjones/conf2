package org.conf2.schema;

/**
 *
 */
public class Rpc extends MetaBase implements MetaCollection, Describable {
    public RpcInput input;
    public RpcOutput output;

    public Rpc(String ident) {
        super(ident);
    }

    @Override
    public Meta getFirstMeta() {
        return input != null ? input : output;
    }

    @Override
    public void addMeta(Meta m) {
        if (m instanceof RpcInput) {
            input = (RpcInput) m;
        } else if (m instanceof RpcOutput) {
            output = (RpcOutput) m;
        } else {
            throw new SchemaError("Can only add input or outputs to rpc");
        }
        if (input != null && output != null) {
            input.setSibling(output);
        }
    }
}

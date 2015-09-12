package org.conf2.schema;

/**
 *
 */
public class Choice extends CollectionBase implements Describable {
    public Choice(String ident) {
        super(ident);
    }

    public ChoiceCase getCase(String caseIdent) {
        return (ChoiceCase) MetaUtil.findByIdent(this, caseIdent);
    }
}

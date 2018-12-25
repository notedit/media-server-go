/*
 * The contents of this file are subject to the Mozilla Public
 * License Version 1.1 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of
 * the License at http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS
 * IS" basis, WITHOUT WARRANTY OF ANY KIND, either express or
 * implied. See the License for the specific language governing
 * rights and limitations under the License.
 *
 * The Original Code is MPEG4IP.
 *
 * The Initial Developer of the Original Code is Cisco Systems Inc.
 * Portions created by Cisco Systems Inc. are
 * Copyright (C) Cisco Systems Inc. 2004.  All Rights Reserved.
 *
 * Contributor(s):
 *      Bill May wmay@cisco.com
 */

#include "src/impl.h"

namespace mp4v2 {
namespace impl {

///////////////////////////////////////////////////////////////////////////////

MP4VpcCAtom::MP4VpcCAtom(MP4File &file)
        : MP4FullAtom(file, "vpcC")
{
    
    AddProperty( new MP4Integer8Property(*this,"profile")); /* 0 */
    AddProperty( new MP4Integer8Property(*this,"level")); /* 1 */
    AddProperty( new MP4BitfieldProperty(*this,"bitDepth" ,4)); /* 2 */
    AddProperty( new MP4BitfieldProperty(*this,"colorSpace", 4)); /* 3 */
    AddProperty( new MP4BitfieldProperty(*this,"chromaSubsampling", 4)); /* 4 */
    AddProperty( new MP4BitfieldProperty(*this,"transferFunction", 3)); /* 5 */
    AddProperty( new MP4BitfieldProperty(*this,"videoFullRangeFlag", 1)); /* 6 */
    AddProperty( new MP4Integer16Property(*this,"codecIntializationDataSize")); /* 7 */
    AddProperty( new MP4BytesProperty(*this,"codecIntializationData",0)); /* 8 */
}

void MP4VpcCAtom::Generate()
{
    MP4FullAtom::Generate();
    ((MP4Integer8Property*)m_pProperties[0])->SetValue(0);
    ((MP4Integer8Property*)m_pProperties[1])->SetValue(0);
    ((MP4BitfieldProperty*)m_pProperties[2])->SetValue(0);
    ((MP4BitfieldProperty*)m_pProperties[3])->SetValue(0);
    ((MP4BitfieldProperty*)m_pProperties[4])->SetValue(0);
    ((MP4BitfieldProperty*)m_pProperties[5])->SetValue(0);
    ((MP4BitfieldProperty*)m_pProperties[6])->SetValue(0);
    ((MP4Integer16Property*)m_pProperties[7])->SetValue(0);
    //((MP4BytesProperty*)m_pProperties[8])->SetValueSize(0,0);
}
///////////////////////////////////////////////////////////////////////////////

}
} // namespace mp4v2::impl

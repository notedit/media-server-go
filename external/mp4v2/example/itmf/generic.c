/* This is an example of iTMF Generic API.
 * WARNING: this program will change/destroy certain tags in an mp4 file.
 */

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <mp4v2/mp4v2.h>

int main( int argc, char** argv )
{
    if( argc != 2 ) {
        printf( "usage: %s file.mp4\n", argv[0] );
        return 1;
    }

    /* open file for modification */
    MP4FileHandle file = MP4Modify( argv[1], MP4_DETAILS_ERROR, 0 );
    if( file == MP4_INVALID_FILE_HANDLE ) {
        printf( "MP4Modify failed\n" );
        return 1;
    }

    /* show existing iTMF items */
    MP4ItmfItemList* list = MP4ItmfGetItems( file );
    printf( "list=%p\n", list );
    if( list ) {
        printf( "list size=%u\n", list->size );
        uint32_t i;
        for( i = 0; i < list->size; i++ ) {
            MP4ItmfItem* item = &list->elements[i];
            printf( "item[%u] type=%s\n", i, item->code );

            if( item->mean )
                printf( "    mean=%s\n", item->mean );
            if( item->name )
                printf( "    name=%s\n", item->name );

            int j;
            for( j = 0; j < item->dataList.size; j++ ) {
                MP4ItmfData* data = &item->dataList.elements[j];
                printf( "    data[%u] typeCode=%u valueSize=%u\n", j, data->typeCode, data->valueSize );
            }
        }

        /* caller responsiblity to free */
        MP4ItmfItemListFree( list );
    }

    /* add bogus item to file */
    {
        /* allocate item with 1 data element */
        MP4ItmfItem* bogus = MP4ItmfItemAlloc( "bogu", 1 );

        const char* const hello = "hello one";

        MP4ItmfData* data = &bogus->dataList.elements[0];
        data->typeCode = MP4_ITMF_BT_UTF8;
        data->valueSize = strlen( hello );
        data->value = (uint8_t*)malloc( data->valueSize );
        memcpy( data->value, hello, data->valueSize );

        /* add to mp4 file */
        MP4ItmfAddItem( file, bogus );

        /* caller responsibility to free */
        MP4ItmfItemFree( bogus );
    }

    /* add bogus item with meaning and name to file */
    {
        /* allocate item with 1 data element */
        MP4ItmfItem* bogus = MP4ItmfItemAlloc( "----", 1 );
        bogus->mean = strdup( "com.garden.Tomato" );
        bogus->name = strdup( "weight" );

        const char* const hello = "hello two";

        MP4ItmfData* data = &bogus->dataList.elements[0];
        data->typeCode = MP4_ITMF_BT_UTF8;
        data->valueSize = strlen( hello );
        data->value = (uint8_t*)malloc( data->valueSize );
        memcpy( data->value, hello, data->valueSize );

        /* add to mp4 file */
        MP4ItmfAddItem( file, bogus );

        /* caller responsibility to free */
        MP4ItmfItemFree( bogus );
    }

    /* free memory associated with structure and close */
    MP4Close( file );

    return 0;
}

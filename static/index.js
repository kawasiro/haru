var GalleryEntryList = React.createClass({
  getInitialState: function() {
    return {
      entries: [],
      opts: {
        gallery: 'hitomi',
        language: 'korean',
        category: '',
        value: '',
        page: 1
      }
    };
  },
  loadEntriesFromServer: function() {
    $.ajax({
      url: this.props.url,
      data: this.state.opts,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({entries: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(url, status, err.toString());
        this.setState({entries: []});
      }.bind(this)
    });
  },
  componentDidMount: function() {
    this.loadEntriesFromServer();
  },
  handleCategoryChange: function(category, value) {
    var opts = this.state.opts;
    opts.category = category;
    opts.value = value;
    opts.page = 1;
    this.setState({opts: opts, entries: []});

    this.loadEntriesFromServer();
  },
  handleLanguageChange: function(language) {
    var opts = this.state.opts
    opts.language = language;
    opts.page = 1;
    this.setState({opts: opts, entries: []});

    this.loadEntriesFromServer();
  },
  handlePageChange: function(page) {
    var page = parseInt(page, 10);
    if(isNaN(page)) {
      page = 1;
    }
    if(page <= 0) {
      page = 1;
    }

    var opts = this.state.opts;
    opts.page = page;
    this.setState({page: page, entries: []});

    this.loadEntriesFromServer();
  },
  handleNextPage: function(e) {
    e.preventDefault();
    var page = this.state.opts.page + 1;
    this.handlePageChange(page);
  },
  handlePrevPage: function(e) {
    e.preventDefault();
    var page = this.state.opts.page - 1;
    this.handlePageChange(page);
  },
  handleReset: function(e) {
    e.preventDefault();
    this.setState(this.getInitialState());
    this.loadEntriesFromServer();
  },
  handleResetPage: function(e) {
    e.preventDefault();
    this.handlePageChange(1);
  },


  render: function() {
    var self = this;
    var entryNodes = this.state.entries.map(function(entry) {
      return (
        <GalleryEntry entry={entry} onCategoryChange={self.handleCategoryChange}/>
      )
    });

    return (
      <div className="galleryEntryList">
      <GalleryListParams opts={this.state.opts}/>
      <button className="btn btn-primary" onClick={this.handleReset}>ResetAll</button>
      <button className="btn btn-default" onClick={this.handlePrevPage}>Prev</button>
      <button className="btn btn-default" onClick={this.handleNextPage}>Next</button>
      <button className="btn btn-default" onClick={this.handlePage}>Rst Page</button>

      <ul>
      {entryNodes}
      </ul>
      </div>
    );
  }
});

var GalleryBox = React.createClass({
  getInitialState: function() {
    return {
      metadata: {},
      imageLinks: [],
      cached: false,
    }
  },
  getId: function() {
    var arr = document.URL.match(/id=([0-9]+)/)
    return arr[1];
  },
  getDetailUrl: function() {
    return '/api/detail/hitomi/' + this.getId();
  },
  getDownloadUrl: function() {
    return '/api/download/hitomi/' + this.getId();
  },
  getEnqueueUrl: function() {
    return '/api/enqueue/hitomi/' + this.getId();
  },
  getReaderUrl: function() {
    return `https://hitomi.la/reader/${this.getId()}.html`
  },
  loadInfoFromServer: function() {
    $.ajax({
      url: this.getDetailUrl(),
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState(data);
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(url, status, err.toString());
        this.setState(this.getInitialState());
      }.bind(this)
    });
  },
  componentDidMount: function() {
    this.loadInfoFromServer();
  },

  render: function() {
    return (
      <div className="galleryBox">
      <GalleryEntry entry={this.state.metadata} onCategoryChange={null}/>
      <div>
      <a className="btn btn-default" href={this.getDownloadUrl()}>Download</a>
      <a className="btn btn-default" href={this.getReaderUrl()}>Read</a>
      <br/>
      <a className="btn btn-default" target="_blank" href={this.getEnqueueUrl()}>
      Enqueue - Cached: {this.state.cached.toString()}
      </a>
      </div>
      </div>
    )
  }
});

var GalleryListParams = React.createClass({
  render: function() {
    var elements = [];
    for(var k in this.props.opts) {
      elements.push(<dt>{k}</dt>);
      elements.push(<dd>{this.props.opts[k]}</dd>);
    }

    return (
      <dl className="dl-horizontal">
      {elements}
      </dl>
    )
  }
});

var GalleryCover = React.createClass({
  url: function() {
    return "/api/proxy/?url=" + this.props.url;
  },
  render: function() {
    return (
      <img src={this.url()}/>
    )
  }
});

var GalleryEntryLabel = React.createClass({
  handleClick: function(e) {
    e.preventDefault();
    this.props.onCategoryChange(this.props.category, this.props.value);
  },
  render: function() {
    return (
      <span>
      <a href="#" onClick={this.handleClick}>
      {this.props.value}
      </a>
      <span> | </span>
      </span>
    )
  }
});

var GalleryEntryTitle = React.createClass({
  url: function() {
    var path = "/detail.html";
    return `${path}?id=${this.props.id}`;
  },
  render: function() {
    return (
      <h2><a href={this.url()} target="_blank">{this.props.title}</a></h2>
    )
  }
});

var GalleryEntry = React.createClass({
  render: function() {
    var onCategoryChange = this.props.onCategoryChange;
    var artistNode = (this.props.entry.artists || []).map(function(v) {
      return (
        <GalleryEntryLabel category={'artist'} value={v}
        onCategoryChange={onCategoryChange}/>
      )
    });

    var tagNode = (this.props.entry.tags || []).map(function(v) {
      return (
        <GalleryEntryLabel category={'tag'} value={v}
        onCategoryChange={onCategoryChange}/>
      )
    });

    var seriesNode = (this.props.entry.series || []).map(function(v) {
      return (
        <GalleryEntryLabel category={'series'} value={v}
        onCategoryChange={onCategoryChange}/>
      )
    });

    var coverNodes = (this.props.entry.covers || []).map(function(url) {
      return <GalleryCover url={url}/>
    });

    return (
      <div className="galleryEntry">
      <GalleryEntryTitle id={this.props.entry.id} title={this.props.entry.title}
      onCategoryChange={onCategoryChange}/>
      <dl className="dl-horizontal">
      <dt>Artist</dt>
      <dd>{artistNode}</dd>
      <dt>Series</dt>
      <dd>{seriesNode}</dd>
      <dt>Type</dt>
      <dd>
      <GalleryEntryLabel category={'type'} value={this.props.entry.type}
      onCategoryChange={onCategoryChange}/>
      </dd>
      <dt>Language</dt>
      <dd>
      <GalleryEntryLabel category={'type'} value={this.props.entry.language}
      onCategoryChange={onCategoryChange}/>
      </dd>
      <dt>Tags</dt>
      <dd>{tagNode}</dd>
      <dt>Date</dt>
      <dd>{this.props.entry.date}</dd>
      </dl>
      <div>{coverNodes}</div>
      </div>
    );
  }
});

if(document.getElementById("content-list") !== null) {
  ReactDOM.render(
    <GalleryEntryList url={'/api/list/hitomi/'}/>,
    document.getElementById('content-list')
  );

} else if(document.getElementById("content-detail") !== null) {
  ReactDOM.render(
    <GalleryBox/>,
    document.getElementById('content-detail')
  );
}
